package mylist

import "testing"

func TestList(t *testing.T) {
	list := New[string]()
	list.Append("jin gui hua ")
	list.Insert(0, "luis")
	list.Insert(1, "is")
	list.Append("hello world")
	t.Log("list :", list)
	if err := list.Move(1, 2); err != nil {
		t.Error(err)
	}
	t.Log("list move 0 to 3:", list)
}
