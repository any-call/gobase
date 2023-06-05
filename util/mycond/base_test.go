package mycond

import "testing"

type MyStruct struct {
	a []int
	b int
}

func TestDeepEQ(t *testing.T) {

	t1 := MyStruct{
		a: []int{1, 3},
		b: 3,
	}

	t2 := MyStruct{
		a: []int{1, 3},
		b: 3,
	}

	t.Log("deep eq :", DeepEQ(t1, &t2))
}
