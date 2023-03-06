package mymap

import "testing"

func TestMap_Insert(t *testing.T) {
	mymap := NewMultiMap[int, string]()
	mymap.Insert(1, "jin")
	mymap.Insert(1, "gui")
	mymap.Insert(1, "hua")
	mymap.Range(func(k int, v string) {
		t.Log(k, v)
	})

	mymap.Remove(1)
	mymap.Remove(1)
	mymap.Remove(1)
	t.Log("after remove ")
	mymap.Range(func(k int, v string) {
		t.Log(k, v)
	})

	t.Log("len:", mymap.Len())
}
