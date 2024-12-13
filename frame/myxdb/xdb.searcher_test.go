package myxdb

import (
	"os"
	"sync"
	"testing"
)

func TestNewXDBFile(t *testing.T) {
	ipSearcher, err := NewXDBSearcher(os.Getenv("SRCFILE"))
	if err != nil {
		t.Error(err)
		return
	}

	var list []string
	for i := 0; i < 100; i++ {
		list = append(list, "43.227.112.128")
	}

	w := &sync.WaitGroup{}
	w.Add(len(list))
	for idx, _ := range list {
		go func(i int, addr string) {
			defer w.Done()
			ret := ipSearcher.Search(list[i])
			t.Logf("%d:%s\n", i, ret)
		}(idx, list[idx])

	}

	w.Wait()

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
