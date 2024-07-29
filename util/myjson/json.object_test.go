package myjson

import (
	"encoding/json"
	"testing"
)

func TestJsonStruct_Scan(t *testing.T) {
	type Tmp struct {
		BaseUrl   string `json:"base_url"`
		ApiKey    string `json:"api_key"`    //基础费
		ApiSecret string `json:"api_secret"` //提供手续费
	}

	aa := Tmp{
		BaseUrl:   "1212",
		ApiKey:    "1212",
		ApiSecret: "12122",
	}
	jb, _ := json.Marshal(aa)
	t.Log("jb is :", string(jb))

	var jObject JsonObject[Tmp]
	if err := jObject.Scan(jb); err != nil {
		t.Error("scan err :", err)
		return
	}

	t.Log("scan is :", jObject.Object)
	jdriver, err := jObject.Value()
	if err != nil {
		t.Error("value err", err)
		return
	}

	t.Log("jDriver is:", jdriver)
}
