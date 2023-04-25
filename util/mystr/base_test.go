package mystr

import "testing"

func TestRemoveSpec(t *testing.T) {
	str := " long long ago ,one morning ,i see "
	ret := RemoveSpec(str, "go")
	t.Log("ret is :", ret)
}

func TestSplitWithRuneLen(t *testing.T) {
	list := SplitWithRuneLen("金贵华12金贵华34金贵华567金贵华8金贵华90", 49)
	for i, _ := range list {
		t.Logf("list[%d]:%s \r\n", i, list[i])
	}
}
