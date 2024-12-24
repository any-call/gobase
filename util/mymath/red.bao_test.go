package mymath

import (
	"testing"
)

func TestHBMath(t *testing.T) {
	list, err := HBMath(100.0, 3, 20)
	if err != nil {
		t.Error(err)
		return
	}

	var total float64
	for i, _ := range list {
		total += list[i]
		t.Logf("%d: %.2f \n", i+1, list[i])
	}

	t.Log("total is :", total)
}
