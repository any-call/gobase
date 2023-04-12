package myvalidator

import (
	"fmt"
	"github.com/any-call/gobase/util/mystr"
	"reflect"
	"regexp"
)

const (
	keyValidTag = "validate"
)

const (
	keyAliasValid = "valid"
)

var (
	reMatchFun *regexp.Regexp
)

func init() {
	reMatchFun = regexp.MustCompile("[a-z]+\\({1}[^()]+\\){1}")
}

func Validate(obj any) error {
	//检测是不是结构体
	myType := reflect.TypeOf(obj)
	myValue := reflect.ValueOf(obj)
	k := myType.Kind()
	if k != reflect.Struct {
		return nil
	}

	//首先检测是不是有 验证标签：validate
	totalField := myType.NumField()
	//fieldList := make([]reflect.StructField, 0)
	for i := 0; i < totalField; i++ {
		//只取必须验证的field
		if vStr, ok := myType.Field(i).Tag.Lookup(keyValidTag); ok {
			if err := validFun(myValue.FieldByName(myType.Field(i).Name), vStr); err != nil {
				return err
			}
		}
	}

	return nil
}

func validFun(val reflect.Value, varStr string) error {
	fmt.Println("input :", varStr, " val:", val)
	list := reMatchFun.FindAllString(varStr, -1)
	for i, _ := range list {
		tmpList := mystr.Split(list[i], []string{"(", ")", ","})
		if len(tmpList) > 1 {
			if tmpList[0] == keyAliasValid {
				//此处要递归 向下验证
			} else {
				validInfo := &ValidInfo{
					val:   val,
					name:  tmpList[0],
					param: tmpList[1:],
				}
				if err := validInfo.Valid(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
