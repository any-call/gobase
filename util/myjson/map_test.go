package myjson

import "testing"

func Test_orderMap(t *testing.T) {
	json := `{
		 "id": "",
    "title": "",
    "titleSub": "",
    "author": "",
    "status": "",
    "publishTime": "",
    "publishTimeStart": "",
    "publishTimeEnd": "",
    "fileUrl": "",
    "orderBy": 100,
    "createTime": "",
    "updateTime": "",
    "createrId": "",
    "updateId": ""
	}`

	list, err := ToOrderMap(json)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range list {
		t.Log(" key:", v.Key, " value:", v.Value)
	}
}
