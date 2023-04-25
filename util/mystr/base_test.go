package mystr

import "testing"

func TestRemoveSpec(t *testing.T) {
	str := " long long ago ,one morning ,i see "
	ret := RemoveSpec(str, "go")
	t.Log("ret is :", ret)
}

func TestSplitWithRuneLen(t *testing.T) {
	list1 := SplitRuneByLen("金贵华12金贵华34金贵华567金贵华8金贵华90", 49)
	for i, _ := range list1 {
		t.Logf("run list[%d]:%s \r\n", i, list1[i])
	}

	list2 := SplitByLen("金贵华", 3)
	for i, _ := range list2 {
		t.Logf(" s list[%d]:%s \r\n", i, list2[i])
	}
}
