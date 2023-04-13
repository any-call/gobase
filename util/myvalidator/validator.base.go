package myvalidator

import (
	"github.com/any-call/gobase/util/mystr"
	"reflect"
	"regexp"
)

const (
	keyValidTag   = "validate"
	keyAliasValid = "valid"
)

var (
	reMatchFun *regexp.Regexp
)

func init() {
	reMatchFun = regexp.MustCompile("[a-z]+\\({1}[^()]+\\){1}")
}

func Validate(obj any) error {
	return check(reflect.TypeOf(obj), reflect.ValueOf(obj))
}

func check(myType reflect.Type, myValue reflect.Value) error {
	//fmt.Printf("enter check type:%v,value:%v\n", myType, myValue)
	//针对指针类 ，转换成 实体
	if myType.Kind() == reflect.Pointer {
		if myValue.IsNil() == false {
			myType = myType.Elem()
			myValue = myValue.Elem()
		} else {
			return nil
		}
	}

	//检测是不是结构体
	if myType.Kind() != reflect.Struct {
		return nil
	}

	//首先检测是不是有 验证标签：validate
	totalField := myType.NumField()
	for i := 0; i < totalField; i++ {
		//只取必须验证的field
		if vStr, ok := myType.Field(i).Tag.Lookup(keyValidTag); ok {
			//正则式取函数相关信息
			listFunInfo := reMatchFun.FindAllString(vStr, -1)
			for j, _ := range listFunInfo {
				tmpList := mystr.Split(listFunInfo[j], []string{"(", ")", ","})
				if len(tmpList) > 1 {
					if tmpList[0] == keyAliasValid {
						//此处要递归 向下验证
						if err := check(myType.Field(i).Type, myValue.FieldByName(myType.Field(i).Name)); err != nil {
							return err
						}
					} else {
						validInfo := &ValidInfo{
							val:   myValue.FieldByName(myType.Field(i).Name),
							name:  tmpList[0],
							param: tmpList[1:],
						}
						if err := validInfo.Valid(); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
