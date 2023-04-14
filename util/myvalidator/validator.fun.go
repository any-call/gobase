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
	switch strings.TrimSpace(strings.ToLower(v.name)) {
	case "min":
		return min(v.val, v.param)

	case "max":
		return max(v.val, v.param)

	case "range":
		return rangeValue(v.val, v.param)

	case "minlength", "minlen":
		return minlength(v.val, v.param)

	case "maxlength", "maxlen":
		return maxlength(v.val, v.param)

	case "arr_minlength", "arr_minlen":
		return arrayMinLen(v.val, v.param)

	case "arr_maxlength", "arr_maxlen":
		return arrayMaxLen(v.val, v.param)

	case "arr_rangelength", "arr_rangelen":
		return arrayRangeLen(v.val, v.param)

	case "map_minlength", "map_minlen":
		return mapMinLen(v.val, v.param)

	case "map_maxlength", "map_maxlen":
		return mapMaxLen(v.val, v.param)

	case "map_rangelength", "map_rangelen":
		return mapRangeLen(v.val, v.param)

	case "rangelength", "rangelen":
		return rangeLength(v.val, v.param)

	case "enum":
		return enum(v.val, v.param)

	case "email", "mail":
		return email(v.val, v.param)

	case "ip":
		return ip(v.val, v.param)

	case "ip4", "ipv4":
		return ip4(v.val, v.param)

	case "ip6", "ipv6":
		return ip6(v.val, v.param)

	case "phone", "mobilephone":
		return phone(v.val, v.param)
	}

	return nil
}

func (v *ValidInfo) String() string {
	return fmt.Sprintf("%v(%v)", v.name, strings.Join(v.param, ","))
}

// 数值类
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
func rangeValue(val reflect.Value, param []string) error {
	if param == nil && len(param) < 2 {
		return nil
	}

	baseStr1 := param[0]
	baseStr2 := param[1]
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			baseVal1, err := myconv.StrToNum[int64](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[int64](baseStr2)
			if err != nil {
				return err
			}

			if val.Int() < baseVal1 || val.Int() > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			baseVal1, err := myconv.StrToNum[uint64](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[uint64](baseStr2)
			if err != nil {
				return err
			}

			if val.Uint() < baseVal1 || val.Uint() > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break

	case reflect.Float32, reflect.Float64:
		{
			baseVal1, err := myconv.StrToNum[float64](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[float64](baseStr2)
			if err != nil {
				return err
			}

			if val.Float() < baseVal1 || val.Float() > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break
	}

	return nil
}

// 字符类
func minlength(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.String:
		{
			if baseVal, err := myconv.StrToNum[int](baseStr); err != nil {
				return err
			} else {
				if len(val.String()) < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func maxlength(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.String:
		{
			if baseVal, err := myconv.StrToNum[int](baseStr); err != nil {
				return err
			} else {
				if len(val.String()) > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func rangeLength(val reflect.Value, param []string) error {
	if param == nil && len(param) < 2 {
		return nil
	}

	baseStr1 := param[0]
	baseStr2 := param[1]
	switch val.Kind() {
	case reflect.String:
		{
			baseVal1, err := myconv.StrToNum[int](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[int](baseStr2)
			if err != nil {
				return err
			}

			if len(val.String()) < baseVal1 || len(val.String()) > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break
	}

	return nil
}

// 枚举值
func enum(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	listEnumItem := strings.Split(baseStr, "|")
	switch val.Kind() {
	case reflect.String:
		{
			mapItem := make(map[string]bool, 10)
			for i, _ := range listEnumItem {
				mapItem[listEnumItem[i]] = true
			}

			if mapItem[val.String()] == false {
				return errors.New(strings.Join(param[1:], " "))
			}
		}
		break

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			mapItem := make(map[int64]bool, 10)
			for i, _ := range listEnumItem {
				if v, err := myconv.StrToNum[int64](listEnumItem[i]); err == nil {
					mapItem[v] = true
				}
			}

			if mapItem[val.Int()] == false {
				return errors.New(strings.Join(param[1:], " "))
			}
		}
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			mapItem := make(map[uint64]bool, 10)
			for i, _ := range listEnumItem {
				if v, err := myconv.StrToNum[uint64](listEnumItem[i]); err == nil {
					mapItem[v] = true
				}
			}

			if mapItem[val.Uint()] == false {
				return errors.New(strings.Join(param[1:], " "))
			}
		}
		break

	case reflect.Float32, reflect.Float64:
		{
			mapItem := make(map[float64]bool, 10)
			for i, _ := range listEnumItem {
				if v, err := myconv.StrToNum[float64](listEnumItem[i]); err == nil {
					mapItem[v] = true
				}
			}

			if mapItem[val.Float()] == false {
				return errors.New(strings.Join(param[1:], " "))
			}
		}
		break
	}

	return nil
}

// 数组类
func arrayMinLen(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		{
			if baseVal, err := myconv.StrToNum[int](baseStr); err != nil {
				return err
			} else {
				if val.Len() < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func arrayMaxLen(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		{
			if baseVal, err := myconv.StrToNum[int64](baseStr); err != nil {
				return err
			} else {
				if int64(val.Len()) > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func arrayRangeLen(val reflect.Value, param []string) error {
	if param == nil && len(param) < 2 {
		return nil
	}

	baseStr1 := param[0]
	baseStr2 := param[1]
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		{
			baseVal1, err := myconv.StrToNum[int64](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[int64](baseStr2)
			if err != nil {
				return err
			}

			if int64(val.Len()) < baseVal1 || int64(val.Len()) > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break
	}

	return nil
}

// map类
func mapMinLen(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Map:
		{
			if baseVal, err := myconv.StrToNum[int](baseStr); err != nil {
				return err
			} else {
				if val.Len() < baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func mapMaxLen(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	baseStr := param[0]
	switch val.Kind() {
	case reflect.Map:
		{
			if baseVal, err := myconv.StrToNum[int64](baseStr); err != nil {
				return err
			} else {
				if int64(val.Len()) > baseVal {
					return errors.New(strings.Join(param[1:], " "))
				}
			}
		}
		break
	}

	return nil
}
func mapRangeLen(val reflect.Value, param []string) error {
	if param == nil && len(param) < 2 {
		return nil
	}

	baseStr1 := param[0]
	baseStr2 := param[1]
	switch val.Kind() {
	case reflect.Map:
		{
			baseVal1, err := myconv.StrToNum[int64](baseStr1)
			if err != nil {
				return err
			}

			baseVal2, err := myconv.StrToNum[int64](baseStr2)
			if err != nil {
				return err
			}

			if int64(val.Len()) < baseVal1 || int64(val.Len()) > baseVal2 {
				return errors.New(strings.Join(param[2:], " "))
			}
		}
		break
	}

	return nil
}

// other 验证
func email(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		{
			if ValidEmail(val.String()) == false {
				return errors.New(strings.Join(param, " "))
			}
		}
		break
	}

	return nil
}

func ip(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		{
			if ValidIP(val.String()) == false {
				return errors.New(strings.Join(param, " "))
			}
		}
		break
	}

	return nil
}

func ip4(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		{
			if ValidIPV4(val.String()) == false {
				return errors.New(strings.Join(param, " "))
			}
		}
		break
	}

	return nil
}

func ip6(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		{
			if ValidIPV6(val.String()) == false {
				return errors.New(strings.Join(param, " "))
			}
		}
		break
	}

	return nil
}

func phone(val reflect.Value, param []string) error {
	if param == nil && len(param) == 0 {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		{
			if ValidPhone(val.String()) == false {
				return errors.New(strings.Join(param, " "))
			}
		}
		break
	}

	return nil
}
