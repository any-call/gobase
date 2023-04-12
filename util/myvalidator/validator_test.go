package myvalidator

import (
	"testing"
)

type MyReq struct {
	ID   int    `json:"id" validate:"min(10,不正确的ID) max(100, 不正确的ID值)"`
	Name string `json:"name" validate:"valid(T)"`
}

func Test_reflect(t *testing.T) {
	req := MyReq{
		ID:   1001,
		Name: "this is test",
	}

	if err := Validate(req); err != nil {
		t.Error(err)
		return
	}

	t.Log("validate: ok ")
}
