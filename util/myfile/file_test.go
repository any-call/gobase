package myfile

import "testing"

func TestReadLines(t *testing.T) {
	file := "/Users/apple/Desktop/广告位分布.txt"
	list, total, err := ReadLines(file)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("list:", len(list), list)
	t.Log("total:", total)
}
