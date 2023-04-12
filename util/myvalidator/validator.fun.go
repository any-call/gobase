package myvalidator

import (
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/myconv"
	"reflect"
	"strings"
)

type ValidInfo struct {
	val   reflect.Value
	name  string   //函数名
	param []string //函数参数
}

func (v *ValidInfo) Valid() error {
	switch v.name {
	case "min":
		return min(v.val, v.param)

	case "max":
		return max(v.val, v.param)
	}

	return nil
}

func (v *ValidInfo) String() string {
	return fmt.Sprintf("%v(%v)", v.name, strings.Join(v.param, ","))
}

func min(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			if baseVal, err := myconv.StrToNum[int64](baseStr); err != nil {
				return err
			} else {
				if val.Int() < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			if baseVal, err := myconv.StrToNum[uint64](baseStr); err != nil {
				return err
			} else {
				if val.Uint() < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break

	case reflect.Float32, reflect.Float64:
		{
			if baseVal, err := myconv.StrToNum[float64](baseStr); err != nil {
				return err
			} else {
				if val.Float() < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}

func max(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			if baseVal, err := myconv.StrToNum[int64](baseStr); err != nil {
				return err
			} else {
				if val.Int() > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			if baseVal, err := myconv.StrToNum[uint64](baseStr); err != nil {
				return err
			} else {
				if val.Uint() > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break

	case reflect.Float32, reflect.Float64:
		{
			if baseVal, err := myconv.StrToNum[float64](baseStr); err != nil {
				return err
			} else {
				if val.Float() > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
