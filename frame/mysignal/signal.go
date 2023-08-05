package mysignal

import (
	"fmt"
	"reflect"
	"sync"
)

func NewSignal() Signal {
	return &signal{}
}

type con struct {
	slot reflect.Value
	cid  int
}

// 定义 Signal结构
type signal struct {
	sync.RWMutex
	cons  []con //标识该信号联接的槽
	conID int   //记录每次联接的增量值
}

// 同步发射 Signal
func (self *signal) Emit(args ...any) {
	self.RLock()
	defer self.RUnlock()

	if self == nil {
		return
	}

	if self.cons == nil || len(self.cons) == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(self.cons))
	for i, _ := range self.cons {
		go func(info con) {
			defer func() {
				wg.Done()
				if p := recover(); p != nil {
					fmt.Println("panic:", p)
				}
			}()

			if paramValues, err := self.setupPublish(info.slot, args...); err != nil {
				return
			} else {
				info.slot.Call(paramValues)
			}
		}(self.cons[i])
	}
	wg.Wait()
	return
}

func (self *signal) Connect(fn any) (int, error) {
	self.Lock()
	defer self.Unlock()

	if self == nil {
		return 0, fmt.Errorf("signal is nil")
	}

	if fn == nil {
		return 0, fmt.Errorf("fn is nil")
	}

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return 0, fmt.Errorf("fn isn't a function ")
	}

	self.conID++
	newConntion := con{
		slot: reflect.ValueOf(fn),
		cid:  self.conID,
	}

	self.cons = append(self.cons, newConntion)
	return self.conID, nil
}

func (self *signal) DisConnect(cid int) error {
	self.Lock()
	defer self.Unlock()

	if self == nil {
		return fmt.Errorf("signal is nil")
	}

	if cid <= 0 || self.cons == nil || len(self.cons) == 0 {
		return fmt.Errorf("incorrect cid:%d", cid)
	}

	self.Lock()
	defer self.Unlock()

	for i, _ := range self.cons {
		if self.cons[i].cid == cid {
			self.cons = append(self.cons[:i], self.cons[i+1:]...)
			return nil
		}
	}

	return nil
}

func (self *signal) setupPublish(callback reflect.Value, args ...any) ([]reflect.Value, error) {
	funcType := callback.Type()
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
