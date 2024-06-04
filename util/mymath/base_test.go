package mymath

import "testing"

func TestFloat64ByPrecision(t *testing.T) {
	t.Log(Float64ByPrecision(3.12, 3, true))
	t.Log(Float64ByPrecision(3.12, 3, false))
}
