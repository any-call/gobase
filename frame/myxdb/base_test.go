package myxdb

import (
	"os"
	"testing"
)

func TestNewXDBFile(t *testing.T) {
	ipSearcher, err := NewXDBSearcher(os.Getenv("XDB_FILE"))
	if err != nil {
		t.Error(err)
		return
	}

	ret := ipSearcher.Search("43.227.112.128")
	t.Log("search restult:", ret)
}
