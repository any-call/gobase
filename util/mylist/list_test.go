package mylist

import "testing"

func TestList(t *testing.T) {
	list := New[string]()
	list.Append("jin gui hua ")
	list.Insert(0, "luis")
	list.Insert(1, "is")
	list.Append("hello world")
	t.Log("list :", list)
	if err := list.Move(3, 0); err != nil {
		t.Error(err)
	}
	t.Log("list move 2 to 0:", list)
}
