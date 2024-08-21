package mycache

import (
	"fmt"
	"testing"
	"time"
)

type CacheTest struct {
	ID   int
	Name string
}

func Test_cache(t *testing.T) {
	c := NewCache(time.Minute)
	c.Set("001", CacheTest{
		ID:   1,
		Name: "luis",
	}, 0)

	if err := c.SaveFile("my.file"); err != nil {
		t.Error("save file err:", err)
		return
	}

	d := NewCache(time.Minute)
	if err := d.LoadFile("my.file"); err != nil {
		t.Error("load file err:", err)
		return
	}

	if obj, b := d.Get("001"); b {
		t.Log("001:", obj)
	}

	t.Log("test ok")
}

func Test_cache1(t *testing.T) {
	c := NewCache(time.Minute)
	c.Set("001", CacheTest{
		ID:   1,
		Name: "luis",
	}, time.Second*5)

	if obj, b := c.Get("001"); b {
		fmt.Println(obj.(CacheTest).ID)
		fmt.Println(obj.(CacheTest).Name)
	} else {
		t.Error("get no exist")
	}

	time.Sleep(time.Second * 5)

	obj, b := c.Get("001")
	t.Log("test ok:", obj, b)
}
