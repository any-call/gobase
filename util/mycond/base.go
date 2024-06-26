package mycond

import (
	"github.com/any-call/gobase/util/myconv"
	"reflect"
)

type BFun func() bool

func Bool(value any) bool {
	b, _ := myconv.ToBool(value)
	return b
}

// And returns true if both a and b are truthy.
func And[T, U any](a T, b U) bool {
	return Bool(a) && Bool(b)
}

// Or returns false if neither a nor b is truthy.
func Or[T, U any](a T, b U) bool {
	return Bool(a) || Bool(b)
}

func If[T any](f BFun, trueVal, falseVal T) T {
	if f != nil {
		if f() {
			return trueVal
		}
	}
	return falseVal
}

func DeepEQ(a any, b any) bool {
	return reflect.DeepEqual(myconv.DirectObj(a), myconv.DirectObj(b))
}
