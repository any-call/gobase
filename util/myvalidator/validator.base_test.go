package myvalidator

import (
	"testing"
)

type MyDeptID int

type MyReq struct {
	ID     int    `json:"id" validate:"min(10,不正确的ID) max(100, 不正确的ID值)"`
	Name   string `json:"name" validate:"valid(T)"`
	Sex    string `json:"sex" validate:"enum(男|女,错误的性别)"`
	MyDept `validate:"valid(T)"`
}

type MyDept struct {
	DeptID MyDeptID `json:"dept_id" validate:"range(1,10,错误的部门ID)"`
	Name   string   `validate:"rangelen(6,10,名称长度必须是6-10)"`
}

func Test_reflect(t *testing.T) {
	req := MyReq{
		ID:   100,
		Name: "this is test",
		Sex:  "男",
		MyDept: MyDept{
			DeptID: 5,
			Name:   "qjdfdf898989",
		},
	}

	if err := Validate(req); err != nil {
		t.Error("err:", err)
		return
	}

	t.Log("validate: ok ")
}
