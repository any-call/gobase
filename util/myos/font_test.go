package myos

import "testing"

func TestGetFontPath(t *testing.T) {
	ret, err := GetFontPath()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ret is :", ret)
}
