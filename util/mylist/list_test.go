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
	list.ResetByArray([]string{"china", " chinese "})
	t.Log("list1 :", list)

	list.AppendByArray([]string{"luis", "very", "good"})
	t.Log("list2 :", list)

	aa := list.TakeHeadN(20)
	t.Log("list3 :", list, ",aa", aa)
}
