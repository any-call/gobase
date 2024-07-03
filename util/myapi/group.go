package myapi

import (
	"fmt"
	"github.com/any-call/gobase/util/mymap"
)

type apiGroup[TYPE any] struct {
	*mymap.Map[string, ApiInfo[TYPE]]
	groupName string
}

func NewApiGroup[DATA any](chLen uint) ApiManager[DATA] {
	ret := &apiGroup[DATA]{}
	ret.Map = mymap.NewMap[string, ApiInfo[DATA]]()
	return ret
}

func (self *apiGroup[DATA]) SetGroup(group string, op DATA) ApiManager[DATA] {
	self.groupName = group
	return self
}

func (self *apiGroup[DATA]) AddGET(path string, module string, op DATA) ApiManager[DATA] {
	self.Insert(genApiKey("GET", self.groupName+path), ApiInfo[DATA]{
		Method: "GET",
		Path:   self.groupName + path,
		Module: module,
		Type:   op,
	})
	return self
}

func (self *apiGroup[DATA]) AddPOST(path string, module string, op DATA) ApiManager[DATA] {
	self.Insert(genApiKey("POST", self.groupName+path), ApiInfo[DATA]{
		Method: "POST",
		Path:   self.groupName + path,
		Module: module,
		Type:   op,
	})
	return self
}

func (self *apiGroup[DATA]) AddPUT(path string, module string, op DATA) ApiManager[DATA] {
	self.Insert(genApiKey("PUT", self.groupName+path), ApiInfo[DATA]{
		Method: "PUT",
		Path:   self.groupName + path,
		Module: module,
		Type:   op,
	})
	return self
}

func (self *apiGroup[DATA]) AddDELETE(path string, module string, op DATA) ApiManager[DATA] {
	self.Insert(genApiKey("DELETE", self.groupName+path), ApiInfo[DATA]{
		Method: "DELETE",
		Path:   self.groupName + path,
		Module: module,
		Type:   op,
	})
	return self
}

func (self *apiGroup[DATA]) ValueBy(method string, routePath string) (ApiInfo[DATA], bool) {
	return self.Map.Value(genApiKey(method, routePath))
}

func (self *apiGroup[TYPE]) List() []ApiInfo[TYPE] {
	_, list := self.Map.ToArray()
	return list
}

func genApiKey(method string, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}
