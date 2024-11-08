package mygitee

import (
	"os"
	"testing"
)

func TestGiteeDev_ListTags(t *testing.T) {
	list, err := NewDevApi(os.Getenv("TOKEN")).
		ListTags("xingyun2024", "ip_node_service_program", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("list is :", list)
}
