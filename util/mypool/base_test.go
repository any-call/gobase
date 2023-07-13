package mypool

import "testing"

type tmpStruct struct {
	Name string
	Sex  int
}

func TestPool_Add(t *testing.T) {
	pool := NewPool[tmpStruct]()
	pool.SetObj(tmpStruct{
		Name: "111",
		Sex:  11,
	})
	pool.SetObj(tmpStruct{
		Name: "222",
		Sex:  11,
	})
	t.Run("get-1", func(t *testing.T) {
		v, b := pool.GetObj()
		t.Log("v:", v, ";b:", b)
	})

	t.Run("get-2", func(t *testing.T) {
		v, b := pool.GetObj()
		t.Log("v:", v, ";b:", b)
	})
}
