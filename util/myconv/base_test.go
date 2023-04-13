package myconv

import (
	"encoding/json"
	"testing"
)

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

func TestToBool(t *testing.T) {
	var aa json.Number = "true"
	if b, err := ToBool(aa); err != nil {
		t.Error(err)
	} else {
		t.Log("b:", b)
	}
}
