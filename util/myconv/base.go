package myconv

import (
	"errors"
	"fmt"
	"strconv"
)

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

//// ToAnyE converts one type to another and returns an error if occurred.
//func ToAnyE[T any](a any) (T, error) {
//	var t T
//	switch any(t).(type) {
//	case bool:
//		v, err := ToBoolE(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case int:
//		v, err := ToIntE(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case int8:
//		v, err := ToInt8E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case int16:
//		v, err := ToInt16E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case int32:
//		v, err := ToInt32E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case int64:
//		v, err := ToInt64E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case uint:
//		v, err := ToUintE(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case uint8:
//		v, err := ToUint8E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case uint16:
//		v, err := ToUint16E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case uint32:
//		v, err := ToUint32E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case uint64:
//		v, err := ToUint64E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case float32:
//		v, err := ToFloat32E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case float64:
//		v, err := ToFloat64E(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	case string:
//		v, err := ToStringE(a)
//		if err != nil {
//			return t, err
//		}
//		t = any(v).(T)
//	default:
//		return t, fmt.Errorf("the type %T is not supported", t)
//	}
//	return t, nil
//}
