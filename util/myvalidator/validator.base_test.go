package myvalidator

import (
	"testing"
)

type MyDeptID int

type MyReq struct {
	ID    int      `json:"id" validate:"gte(10,不正确的ID)"`
	Name  string   `json:"name" validate:"len(1,用户名长度必须为1)"`
	Sex   string   `json:"sex" validate:"enum(男|女,错误的性别)"`
	Date  string   `json:"date" validate:"date(2006-01-02,无效的日期格式)"`
	MyArr []MyDept `validate:"arr_minlen(1,入参数组不能为空) valid(T)"`
	MyMap map[string]MyDept
}

type MyDept struct {
	DeptID MyDeptID `json:"dept_id" validate:"range(1,10,错误的部门ID)"`
	Name   string   `validate:"rangelen(6,10,名称长度必须是6-10)"`
}

func Test_validator(t *testing.T) {
	req := MyReq{
		ID:   9,
		Name: "1",
		Sex:  "男",
		Date: "1979-12-20",
		MyArr: []MyDept{{
			DeptID: 5,
			Name:   "lu889i",
		}, {
			DeptID: 1,
			Name:   "jinguihua",
		}},
		MyMap: map[string]MyDept{
			"12": {
				DeptID: 1001,
				Name:   "good",
			},
		},
	}

	if err := Validate(req); err != nil {
		t.Error("err:", err)
		return
	}

	t.Log("validate: ok ")
}
