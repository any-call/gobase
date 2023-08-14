package mylist

import (
	"testing"
)

func TestList(t *testing.T) {
	list := NewList[string]()
	list.Append("jin gui hua ")
	list.Insert(0, "luis")
	list.Insert(1, "is")
	list.Append("hello world")
	t.Log("list :", list)
	//if err := list.Move(3, 0); err != nil {
	//	t.Error(err)
	//}

	//list.SwapItemsAt(3, 2)

	list.Range(func(index int, v string) {
		t.Log("index", index, " value:", v)
	})

	t.Log("list :", list)
	//var err error
	//var v string
	//for err == nil {
	//	if v, err = list.TakeAt(0); err == nil {
	//		t.Log("v:", v)
	//	}
	//
	//}
}
