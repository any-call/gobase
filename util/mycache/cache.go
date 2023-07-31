package mycache

import (
	"encoding/gob"
	"fmt"
	"github.com/any-call/gobase/util/mymap"
	"os"
	"time"
)

type cache struct {
	items *mymap.Map[string, Item]
}

func NewCache() Cache {
	return &cache{items: mymap.NewMap[string, Item]()}
}

func (self *cache) Set(k string, v any, d time.Duration) {
	if v == nil {
		return
	}

	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	self.items.Insert(k, Item{
		Object:     v,
		Expiration: e,
	})
}

func (self *cache) Get(k string) (any, bool) {
	obj, b := self.items.Value(k)
	if !b {
		return nil, false
	}

	if obj.Expired() {
		self.items.TakeAt(k)
		return nil, false
	}

	return obj.Object, true
}

func (self *cache) HasKey(k string) bool {
	_, b := self.items.Value(k)
	return b
}

func (self *cache) UpdateExpired(k string, d time.Duration) error {
	v, b := self.items.Value(k)
	if !b {
		return fmt.Errorf("%s is not exist", k)
	}

	self.Set(k, v.Object, d)
	return nil
}

func (self *cache) Flush() {
	keys, _ := self.items.ToArray()
	for i, _ := range keys {
		self.items.Remove(keys[i])
	}
}

func (self *cache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := gob.NewDecoder(fp)
	items := map[string]Item{}
	err = dec.Decode(&items)
	if err == nil {
		for k, v := range items {
			self.items.Insert(k, v)
		}
	}
	return err
}

func (self *cache) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(fp)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()

	tmp := make(map[string]Item)
	self.items.Range(func(key string, value Item) {
		tmp[key] = value
		gob.Register(value.Object)
	})

	err = enc.Encode(&tmp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}
