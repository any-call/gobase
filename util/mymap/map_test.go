package mymap

import (
	"fmt"
	"strings"
	"testing"
)

func TestMap_Insert(t *testing.T) {
	mymap := NewMultiMap[int, string]()
	mymap.Insert(1, "jin")
	mymap.Insert(1, "gui")
	mymap.Insert(1, "hua")
	mymap.Insert(1, "is")
	mymap.Insert(1, "very")
	mymap.Insert(1, "good")

	mymap.RemoveAtIndex(1, 6)
	if list, b := mymap.Values(1); b {
		fmt.Println("after remove :", strings.Join(list, " "))
	}

	mymap.Remove(1)
	mymap.Remove(1)
	mymap.Remove(1)
	t.Log("after remove ")
	mymap.Range(func(k int, v string) {
		t.Log(k, v)
	})

	t.Log("len:", mymap.Len())
}
