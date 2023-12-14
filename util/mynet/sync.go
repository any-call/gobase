package mynet

import (
	"fmt"
	"gitee.com/any-call/gobase/util/mymap"
	"reflect"
	"sync"
)

type OPInfo struct {
	Key      string
	Fn       any
	Args     []any
	callback reflect.Value
	cbArgs   []reflect.Value
}

func (self *OPInfo) Check() error {
	if self.Fn == nil {
		return fmt.Errorf("fun is nil")
	}

	if reflect.TypeOf(self.Fn).Kind() != reflect.Func {
		return fmt.Errorf("callback isn't a function ")
	}
	self.callback = reflect.ValueOf(self.Fn)
	funcType := self.callback.Type()
	inCount := funcType.NumIn()
	if inCount > 0 {
		if self.Args == nil {
			return fmt.Errorf("fun %s :incorrect args count ", self.Key)
		}
		if inCount != len(self.Args) {
			return fmt.Errorf("fun input parameter err:%v", self.Args)
		}

		self.cbArgs = make([]reflect.Value, len(self.Args))
		for i, v := range self.Args {
			if v == nil {
				self.cbArgs[i] = reflect.New(funcType.In(i)).Elem()
			} else {
				self.cbArgs[i] = reflect.ValueOf(v)
			}
		}
	}

	return nil
}

func (self *OPInfo) Exec() any {
	list := self.callback.Call(self.cbArgs)
	if list != nil {
		ret := make([]any, len(list))
		for i, _ := range list {
			if list[i].CanInterface() {
				ret[i] = list[i].Interface()
			} else if list[i].CanInt() {
				ret[i] = list[i].Int()
			} else if list[i].CanFloat() {
				ret[i] = list[i].Float()
			} else if list[i].Kind() == reflect.String {
				ret[i] = list[i].String()
			}
		}

		if len(ret) == 1 {
			return ret[0]
		}
		return ret
	}

	return nil
}

func SyncOP(listOP ...OPInfo) map[string]any {
	wg := sync.WaitGroup{}
	wg.Add(len(listOP))

	ret := mymap.NewMap[string, any]()
	for i, _ := range listOP {
		go func(op OPInfo) {
			defer func() {
				wg.Done()
				if p := recover(); p != nil {
					ret.Insert(op.Key, p)
				}
			}()

			if err := op.Check(); err != nil {
				ret.Insert(op.Key, err)
			} else {
				ret.Insert(op.Key, op.Exec())
			}

		}(listOP[i])
	}

	wg.Wait()
	return ret.ToMap()
}

func SyncOPByNum(goNum uint, listOP ...OPInfo) map[string]any {
	if goNum == 0 {
		goNum = 1
	}
	limiter := make(chan struct{}, goNum)

	wg := sync.WaitGroup{}
	wg.Add(len(listOP))

	ret := mymap.NewMap[string, any]()
	for i, _ := range listOP {
		limiter <- struct{}{}
		go func(op OPInfo) {
			defer func() {
				<-limiter
				wg.Done()
				if p := recover(); p != nil {
					ret.Insert(op.Key, p)
				}
			}()

			if err := op.Check(); err != nil {
				ret.Insert(op.Key, err)
			} else {
				ret.Insert(op.Key, op.Exec())
			}

		}(listOP[i])
	}

	wg.Wait()
	return ret.ToMap()
}
