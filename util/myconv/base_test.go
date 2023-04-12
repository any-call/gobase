package myconv

import "testing"

func TestStrToNum(t *testing.T) {
	if v, err := StrToNum[int]("123"); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("str to num :", v)
	}

	if v, err := StrToNum[float32]("adfd"); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("str to num :", v)
	}
}
