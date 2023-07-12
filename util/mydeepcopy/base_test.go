package mydeepcopy

import "testing"

type baseStr struct {
	Name int
	BV   []int
}

func TestDeepCopy(t *testing.T) {
	//slice
	var a []interface{} = []interface{}{"a", 42, true, 4.32}
	b := Copy(a)
	t.Log("slice a:", a)
	t.Log("slice b:", b)

	var mapA map[any]any = make(map[any]any, 10)
	mapA["aa"] = []string{"af", "luis"}
	mapA[121] = "dfdf"
	mapB := Copy(mapA)
	t.Log("map a:", mapA)
	t.Log("map b:", mapB)

	var structA baseStr = baseStr{
		Name: 0,
		BV:   []int{1, 2},
	}

	structB := Copy(structA)
	t.Log("struct a:", &(structA.BV))
	t.Log("struct b:", (structB.(baseStr).BV))

	t.Log("run ok")
}
