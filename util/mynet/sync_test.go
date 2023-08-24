package mynet

import (
	"testing"
)

func Test_syncOP(t *testing.T) {
	fn1 := func(key string) any {
		return []int{1, 2, 3}
	}

	fn2 := func(key string) any {
		return struct {
			ID   int
			Name string
		}{ID: 001, Name: "luis.jin"}
	}

	fn3 := func(key string) any {
		return map[string]int{
			"001": 223,
		}
	}
	fn4 := func() bool {
		return false
	}

	list := make([]OPInfo, 0)
	list = append(list, OPInfo{
		Key:  "001",
		Fn:   fn1,
		Args: []any{"1121"},
	}, OPInfo{
		Key:  "002",
		Fn:   fn2,
		Args: []any{"1121"},
	}, OPInfo{
		Key:  "003",
		Fn:   fn3,
		Args: []any{"1121"},
	}, OPInfo{
		Key:  "004",
		Fn:   fn4,
		Args: nil,
	})

	ret := SyncOP(list...)
	t.Log("run ok...", ret)
	ret = SyncOPByNum(1, list...)
	t.Log("run ok 11...", ret)
}
