package mycache

import "testing"

type CacheTest struct {
	ID   int
	Name string
}

func Test_cache(t *testing.T) {
	c := NewCache()
	c.Set("001", CacheTest{
		ID:   1,
		Name: "luis",
	}, 0)

	if err := c.SaveFile("my.file"); err != nil {
		t.Error("save file err:", err)
		return
	}

	d := NewCache()
	if err := d.LoadFile("my.file"); err != nil {
		t.Error("load file err:", err)
		return
	}

	if obj, b := d.Get("001"); b {
		t.Log("001:", obj)
	}

	t.Log("test ok")
}
