package myconfig

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Encode 将配置值转换为可用于表存储的列表数据。
// 结构体的每个一级导出字段对应一条 Item，复杂字段使用 JSON 保存。
func Encode(value any) ([]Item, error) {
	if value == nil {
		return nil, fmt.Errorf("myconfig: value must not be nil")
	}

	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return []Item{{Key: rootKey, Type: TypeNil}}, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		item, err := encodeItem(rootKey, v)
		if err != nil {
			return nil, err
		}
		return []Item{item}, nil
	}

	t := v.Type()
	items := make([]Item, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		item, err := encodeItem(fieldType.Name, v.Field(i))
		if err != nil {
			return nil, fmt.Errorf("myconfig: key=%s: %w", fieldType.Name, err)
		}
		items = append(items, item)
	}
	return items, nil
}

func encodeItem(key string, value reflect.Value) (Item, error) {
	item := Item{Key: key}

	switch value.Kind() {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if value.IsNil() {
			item.Type = TypeNil
			return item, nil
		}
	}

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		item.Type = TypeString
		item.Value = value.String()
	case reflect.Bool:
		item.Type = TypeBool
		item.Value = strconv.FormatBool(value.Bool())
	case reflect.Int:
		item.Type = TypeInt
		item.Value = strconv.FormatInt(value.Int(), 10)
	case reflect.Int8:
		item.Type = TypeInt8
		item.Value = strconv.FormatInt(value.Int(), 10)
	case reflect.Int16:
		item.Type = TypeInt16
		item.Value = strconv.FormatInt(value.Int(), 10)
	case reflect.Int32:
		item.Type = TypeInt32
		item.Value = strconv.FormatInt(value.Int(), 10)
	case reflect.Int64:
		item.Type = TypeInt64
		item.Value = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		item.Type = TypeUint
		item.Value = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint8:
		item.Type = TypeUint8
		item.Value = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint16:
		item.Type = TypeUint16
		item.Value = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint32:
		item.Type = TypeUint32
		item.Value = strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint64:
		item.Type = TypeUint64
		item.Value = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		item.Type = TypeFloat32
		item.Value = strconv.FormatFloat(value.Float(), 'g', -1, 32)
	case reflect.Float64:
		item.Type = TypeFloat64
		item.Value = strconv.FormatFloat(value.Float(), 'g', -1, 64)
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		data, err := json.Marshal(value.Interface())
		if err != nil {
			return Item{}, err
		}
		item.Type = TypeJSON
		item.Value = string(data)
	default:
		return Item{}, fmt.Errorf("unsupported type %s", value.Type())
	}

	return item, nil
}
