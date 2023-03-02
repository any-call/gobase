package myconv

import "testing"

func TestIPV42Long(t *testing.T) {
	aa := IPV42Long("abcd")
	t.Log("aa:", aa)
}
