package mystr

import "testing"

func TestRemoveSpec(t *testing.T) {
	str := " long long ago ,one morning ,i see "
	ret := RemoveSpec(str, "go")
	t.Log("ret is :", ret)
}
