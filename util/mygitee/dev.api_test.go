package mygitee

import (
	"os"
	"testing"
)

func TestGiteeDev_ListTags(t *testing.T) {
	err := NewDevApi(os.Getenv("TOKEN")).GetZipFile("jinguihua",
		"botApi",
		"v1.0.0",
		"/Users/luisjin/Desktop/temp/1.zip")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("get ok")
}
