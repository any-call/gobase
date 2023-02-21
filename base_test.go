package cmd

import (
	"github/any-call/gobase/util/mylist"
	"testing"
)

// 测试并集
func TestList_union(t *testing.T) {
	a := []string{"1", "3", "5"}
	b := []string{"11", "3", "15"}

	c := mylist.Union[string](a, b)
	t.Log("union:", c)
}

func TestList_intersect(t *testing.T) {
	a := []string{"11", "3", "5", "32"}
	b := []string{"11", "32", "15"}

	c := mylist.Intersect[string](a, b)
	t.Log("intersect:", c)
}

func TestList_difference(t *testing.T) {
	a := []string{"11", "3", "5"}
	b := []string{"11", "3", "15"}

	c := mylist.Difference[string](a, b)
	t.Log("union:", c)
}

func TestList_removeDuplicItem(t *testing.T) {
	a := []string{"11", "11", "12", "12", "13", "14", "15", "15"}
	a1 := mylist.RemoveDuplicItem[string](a)
	t.Log("a:", a)
	t.Log("a1:", a1)
}
