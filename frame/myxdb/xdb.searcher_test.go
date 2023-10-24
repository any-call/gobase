package myxdb

import (
	"os"
	"testing"
)

func TestNewXDBFile(t *testing.T) {
	ipSearcher, err := NewXDBSearcher("bb.xx")
	if err != nil {
		t.Error(err)
		return
	}

	ret := ipSearcher.Search("43.227.112.128")
	t.Log("search restult:", ret)
}

func TestNewXDBMaker(t *testing.T) {
	ret, err := NewXDBMakerByVector(os.Getenv("SRCFILE"), "bb.xx")
	if err != nil {
		t.Error(err)
		return
	}

	if err := ret.GenXDBFile(); err != nil {
		t.Error("gen xdb file err:", err)
		return
	}

	t.Log("new xdb file:bb.xx")
}
