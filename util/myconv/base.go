package myconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 去掉引用，返回对象本身
func DirectObj(a any) any {
	if a == nil {
		return nil
	}

	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		return a
	}

	//说明是指针类 :检测是不是空指针
	v := reflect.ValueOf(a)
	if v.IsNil() {
		return nil
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	return v.Interface()
}

func ToAny[T any](i any) (T, error) {
	var ret T
	myType := reflect.TypeOf(ret)
	switch myType.Kind() {
	case reflect.Int:
		{
			tmp, err := ToNum[int](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Int8:
		{
			tmp, err := ToNum[int8](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Int16:
		{
			tmp, err := ToNum[int16](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Int32:
		{
			tmp, err := ToNum[int32](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Int64:
		{
			tmp, err := ToNum[int64](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Uint:
		{
			tmp, err := ToNum[uint](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Uint8:
		{
			tmp, err := ToNum[uint8](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Uint16:
		{
			tmp, err := ToNum[uint16](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Uint32:
		{
			tmp, err := ToNum[uint32](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Uint64:
		{
			tmp, err := ToNum[uint64](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Float32:
		{
			tmp, err := ToNum[float32](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.Float64:
		{
			tmp, err := ToNum[float64](i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}

	case reflect.String:
		{
			tmp := ToStr(i)
			return any(tmp).(T), nil
		}

	case reflect.Bool:
		{
			tmp, err := ToBool(i)
			if err != nil {
				return ret, err
			}
			return any(tmp).(T), nil
		}
	}

	return ret, fmt.Errorf("unable to cast %#v of type %T to  %T", i, i, ret)
}

func ToBool(i any) (bool, error) {
	i = DirectObj(i)
	switch reflect.ValueOf(i).Kind() {
	case reflect.Bool:
		return i.(bool), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if reflect.ValueOf(i).Int() != 0 {
			return true, nil
		}

		return false, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if reflect.ValueOf(i).Uint() != 0 {
			return true, nil
		}
		return false, nil

	case reflect.Float32, reflect.Float64:
		if reflect.ValueOf(i).Float() != 0 {
			return true, nil
		}
		return false, nil

	case reflect.String:
		return strconv.ParseBool(reflect.ValueOf(i).String())
	}

	return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
}

func ToNum[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](i any) (T, error) {
	i = DirectObj(i)
	return StrToNum[T](fmt.Sprintf("%v", i))
}

func ToStr(i any) string {
	i = DirectObj(i)
	return fmt.Sprintf("%v", i)
}

func StrToNum[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](str string) (ret T, err error) {
	switch any(ret).(type) {
	case int, int8, int16, int32, int64:
		{
			v, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return ret, err
			}

			switch any(ret).(type) {
			case int:
				ret = any(int(v)).(T)
				break

			case int8:
				ret = any(int8(v)).(T)
				break

			case int16:
				ret = any(int16(v)).(T)
				break

			case int32:
				ret = any(int32(v)).(T)
				break

			default:
				ret = any(v).(T)
				break
			}
		}
		break

	case uint, uint16, uint32, uint64:
		{
			v, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return ret, err
			}

			switch any(ret).(type) {
			case uint:
				ret = any(uint(v)).(T)
				break
			case uint16:
				ret = any(uint16(v)).(T)
				break
			case uint32:
				ret = any(uint32(v)).(T)
				break
			default:
				ret = any(v).(T)
				break
			}
		}
		break

	case float32, float64:
		{
			v, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return ret, err
			}

			switch any(ret).(type) {
			case float32:
				ret = any(float32(v)).(T)
				break
			default:
				ret = any(v).(T)
			}

		}
		break

	default:
		return ret, errors.New(fmt.Sprintf("not match type:%v", ret))
	}

	return ret, nil
}
