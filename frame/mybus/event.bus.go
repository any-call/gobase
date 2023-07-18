package mybus

import (
	"fmt"
	"github.com/any-call/gobase/util/mymap"
	"reflect"
	"sync"
)

// event Handler
type eventHandler struct {
	sync.Mutex
	callBack     reflect.Value
	flagOnce     bool //只运行一次
	flagAsync    bool //异步执行
	flagSequence bool //异步序列执行标识
}

type eventBus struct {
	sync.Mutex
	handlers mymap.MultiMap[string, *eventHandler]
	wg       sync.WaitGroup
}

func NewEventBus() EventBus {
	return &eventBus{
		Mutex:    sync.Mutex{},
		handlers: mymap.MultiMap[string, *eventHandler]{},
		wg:       sync.WaitGroup{},
	}
}

// 订阅
func (self *eventBus) doSubscribe(key string, fn any, ehandler *eventHandler) error {
	if fn == nil || ehandler == nil {
		return fmt.Errorf("callback is nil")
	}

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("callback isn't a function ")
	}

	self.handlers.Insert(key, ehandler)
	return nil
}

func (self *eventBus) Subscribe(key string, fn any) error {
	return self.doSubscribe(key, fn, &eventHandler{
		Mutex:        sync.Mutex{},
		callBack:     reflect.ValueOf(fn),
		flagOnce:     false,
		flagAsync:    false,
		flagSequence: false,
	})
}

func (self *eventBus) SubscribeAsync(key string, fn any, sequence bool) error {
	return self.doSubscribe(key, fn, &eventHandler{
		Mutex:        sync.Mutex{},
		callBack:     reflect.ValueOf(fn),
		flagOnce:     false,
		flagAsync:    true,
		flagSequence: sequence,
	})
}

func (self *eventBus) SubscribeOnce(key string, fn any) error {
	return self.doSubscribe(key, fn, &eventHandler{
		Mutex:        sync.Mutex{},
		callBack:     reflect.ValueOf(fn),
		flagOnce:     true,
		flagAsync:    false,
		flagSequence: false,
	})
}

func (self *eventBus) SubscribeOnceAsync(key string, fn any, sequence bool) error {
	return self.doSubscribe(key, fn, &eventHandler{
		Mutex:        sync.Mutex{},
		callBack:     reflect.ValueOf(fn),
		flagOnce:     true,
		flagAsync:    true,
		flagSequence: sequence,
	})
}

func (self *eventBus) Unsubscribe(key string, fn any) {
	self.Lock()
	defer self.Unlock()

	if fn == nil {
		return
	}

	cbFun := reflect.ValueOf(fn)
	var idx int = -1

	self.handlers.SearchKey(key, func(index int, value *eventHandler) bool {
		if value.callBack.Type() == cbFun.Type() &&
			value.callBack.Pointer() == cbFun.Pointer() {
			idx = index
			return true
		}

		return false
	})

	if idx >= 0 {
		self.handlers.RemoveAtIndex(key, idx)
	}
}

func (self *eventBus) Publish(key string, args ...any) {
	self.Lock()
	defer func() {
		self.Unlock()
		p := recover()
		if p != nil {
			fmt.Println("publish panic:", p)
		}
	}()

	if list, b := self.handlers.Values(key); b {
		copyHandlers := make([]*eventHandler, len(list))
		copy(copyHandlers, list)
		for i, handler := range copyHandlers {
			if handler.flagOnce {
				self.handlers.RemoveAtIndex(key, i)
			}

			if !handler.flagAsync {
				self.doPublish(handler, args...)
			} else {
				self.wg.Add(1)
				if handler.flagSequence {
					self.Unlock()
					handler.Lock()
					self.Lock()
				}
				go self.doPublishAsync(handler, args...)
			}
		}
	}
}

func (self *eventBus) HasCallback(key string) bool {
	return self.handlers.HasKey(key)
}

func (self *eventBus) WaitAsync() {
	self.wg.Wait()
}

func (self *eventBus) doPublish(handler *eventHandler, args ...any) error {
	passedArguments, err := self.setupPublish(handler, args...)
	if err != nil {
		return err
	}

	handler.callBack.Call(passedArguments)
	return nil
}

func (self *eventBus) doPublishAsync(handler *eventHandler, args ...any) error {
	defer self.wg.Done()
	if handler.flagSequence {
		defer handler.Unlock()
	}

	return self.doPublish(handler, args...)
}

func (self *eventBus) setupPublish(handler *eventHandler, args ...any) ([]reflect.Value, error) {
	funcType := handler.callBack.Type()
	inCount := funcType.NumIn()
	if inCount != len(args) {
		return nil, fmt.Errorf("fun input parameter err:%v", args)
	}

	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments, nil
}
