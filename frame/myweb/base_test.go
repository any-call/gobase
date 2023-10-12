package myweb

import (
	"testing"
)

func Test_myweb(t *testing.T) {
	type (
		Tmp struct {
			ID   int    `validate:"min(5,invalid value)"`
			Name string `validate:"min_length(5,字符数不足)"`
		}
	)

	aa := Tmp{
		ID:   8,
		Name: "12322",
	}

	err := validate(&aa)
	t.Log("valid err:", err)
}
