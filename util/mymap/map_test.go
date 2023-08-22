package mymap

import (
	"testing"
)

func TestMap_Insert(t *testing.T) {
	mymap := NewMap[int, string]()
	mymap.Insert(1, "jin")
	mymap.Insert(2, "gui")
	mymap.Insert(3, "hua")
	mymap.Insert(4, "is")
	mymap.Insert(5, "very")
	mymap.Insert(6, "good")

	mymap.Range(func(k int, v string) {
		t.Log(k, v)
	})

	tmpMap := make(map[int]string)
	tmpMap[10] = "china"
	tmpMap[11] = "chinese"
	mymap.ResetByMap(tmpMap)

	t.Log("mymap is :")
	mymap.Range(func(k int, v string) {
		t.Log(k, v)
	})
}
