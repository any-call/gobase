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

	ret := SyncOP(OPInfo{
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
	t.Log("run ok...", ret)
}
