package mymap

import (
	"testing"
)

func TestMap_Insert(t *testing.T) {
	mymap := NewMap[int, string]()
	mymap.Reset(len([]int{}))
	mymap.Insert(1, "jin")
	mymap.Insert(2, "gui")
	mymap.Insert(3, "hua")
	mymap.Insert(4, "is")
	mymap.Insert(5, "very")
	mymap.Insert(6, "good")
	
	t.Log("mymap is :", mymap.mapList)
	v, b := mymap.TakeOne()
	t.Log("v ,b", v, b)
	t.Log("mymap take one ", mymap.mapList)
}
