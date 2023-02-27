package cmd

import (
	"github.com/any-call/gobase/util/mylist"
	"github.com/any-call/gobase/util/myvalidator"
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

func Test_ValidFqdn(t *testing.T) {
	b1 := myvalidator.ValidFqdn("baidu.com")
	b2 := myvalidator.ValidFqdn("aa.baidu.com")
	t.Log("b1", b1)
	t.Log("b2", b2)
}

func Test_ValidEmail(t *testing.T) {
	b1 := myvalidator.ValidEmail("baidu.com")
	b2 := myvalidator.ValidEmail("12121212@cccc.com")
	t.Log("b1", b1)
	t.Log("b2", b2)
}
