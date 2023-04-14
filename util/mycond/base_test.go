package mycond

import "testing"

func TestIf(t *testing.T) {
	a := "false"
	b := reflectValue(&a)
	t.Log("b :", b)
}
