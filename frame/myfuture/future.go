package myfuture

import (
	"fmt"
)

//○ wait ： 同时执行多个，异步
//○ then ：结果回调
//○ catchError : 例外
//○ whenComplete : 执行结束
//○ timeout : 延时执行

type future struct {
	onThenCB   func()
	onCatchErr func(err error)
	onComplete func()
}

func Start(f func() error) Future {
	if f == nil {
		panic("empty func")
	}

	fut := &future{}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				//说明函数出了例外
				if fut.onCatchErr != nil {
					fut.onCatchErr(fmt.Errorf("%v", r))
				}

				if fut.onComplete != nil {
					fut.onComplete()
				}
			}
		}()

		retFail := f()
		if retFail != nil {
			if fut.onCatchErr != nil {
				fut.onCatchErr(retFail)
			}
		} else {
			if fut.onThenCB != nil {
				fut.onThenCB()
			}
		}

		if fut.onComplete != nil {
			fut.onComplete()
		}
	}()

	return fut
}

func (self *future) Then(f func()) Future {
	self.onThenCB = f
	return self
}

func (self *future) Catch(f func(error)) Future {
	self.onCatchErr = f
	return self
}

func (self *future) Complete(f func()) Future {
	self.onComplete = f
	return self
}
