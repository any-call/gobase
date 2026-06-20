package myconfig

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Decode 将列表数据还原到目标配置模型。
// out 必须是非 nil 指针；缺失字段保留 out 中原有的值，未知 Key 会被忽略。
func Decode(items []Item, out any) error {
	if out == nil {
		return fmt.Errorf("myconfig: output must not be nil")
	}

	target := reflect.ValueOf(out)
	if target.Kind() != reflect.Ptr || target.IsNil() {
		return fmt.Errorf("myconfig: output must be a non-nil pointer")
	}

	itemMap := make(map[string]Item, len(items))
	for _, item := range items {
		if item.Key == "" {
			return fmt.Errorf("myconfig: item key must not be empty")
		}
		if _, exists := itemMap[item.Key]; exists {
			return fmt.Errorf("myconfig: duplicate key %q", item.Key)
		}
		itemMap[item.Key] = item
	}

	target = target.Elem()
	base := indirectType(target.Type())
	if base.Kind() != reflect.Struct {
		item, ok := itemMap[rootKey]
		if !ok {
			return fmt.Errorf("myconfig: root key %q not found", rootKey)
		}
		if err := decodeItem(item, target); err != nil {
			return fmt.Errorf("myconfig: key=%s: %w", rootKey, err)
		}
		return nil
	}

	for target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		target = target.Elem()
	}

	t := target.Type()
	for i := 0; i < target.NumField(); i++ {
		fieldType := t.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		item, ok := itemMap[fieldType.Name]
		if !ok {
			continue
		}
		if err := decodeItem(item, target.Field(i)); err != nil {
			return fmt.Errorf("myconfig: key=%s: %w", fieldType.Name, err)
		}
	}
	return nil
}

func decodeItem(item Item, target reflect.Value) error {
	if !target.CanSet() {
		return fmt.Errorf("target type %s cannot be set", target.Type())
	}

	if item.Type == TypeNil {
		switch target.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
			target.Set(reflect.Zero(target.Type()))
			return nil
		default:
			return fmt.Errorf("nil cannot be assigned to %s", target.Type())
		}
	}

	if target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		return decodeItem(item, target.Elem())
	}

	if item.Type == TypeJSON {
		if !target.CanAddr() {
			return fmt.Errorf("target type %s is not addressable", target.Type())
		}
		if err := json.Unmarshal([]byte(item.Value), target.Addr().Interface()); err != nil {
			return err
		}
		return nil
	}

	switch target.Kind() {
	case reflect.String:
		target.SetString(item.Value)
	case reflect.Bool:
		value, err := strconv.ParseBool(item.Value)
		if err != nil {
			return err
		}
		target.SetBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(item.Value, 10, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(item.Value, 10, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(item.Value, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetFloat(value)
	default:
		return fmt.Errorf("type %q cannot be decoded into %s", item.Type, target.Type())
	}
	return nil
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
