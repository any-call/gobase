package myos

import "testing"

func TestGetFontPath(t *testing.T) {
	ret, err := GetSystemFontPath()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ret is :", ret)
}
