package myconv

import (
	"encoding/json"
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

func ToBool(i any) (bool, error) {
	i = DirectObj(i)
	switch b := i.(type) {
	case bool:
		return b, nil

	case int, int8, int16, int32, int64:
		if reflect.ValueOf(i).Int() != 0 {
			return true, nil
		}
		return false, nil

	case uint, uint8, uint16, uint32, uint64:
		if reflect.ValueOf(i).Uint() != 0 {
			return true, nil
		}
		return false, nil

	case float32, float64:
		if reflect.ValueOf(i).Float() != 0 {
			return true, nil
		}
		return false, nil

	case string:
		return strconv.ParseBool(i.(string))
	}

	return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
}

func ToInt(i any) (int, error) {
	i = DirectObj(i)

	switch v := i.(type) {
	case int:
		return int(v), nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case json.Number:
		return strconv.Atoi(string(v))
	}

	return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
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
