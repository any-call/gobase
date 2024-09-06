package myfuture

import (
	"fmt"
)

//○ wait ： 同时执行多个，异步
//○ then ：结果回调
//○ catchError : 例外
//○ whenComplete : 执行结束
//○ timeout : 延时执行

type future[T any] struct {
	onThenCB   func(T)
	onCatchErr func(err error)
	onComplete func()
}

func Start[T any](f func() (T, error)) Future[T] {
	if f == nil {
		panic("empty func")
	}

	fut := &future[T]{}
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

		retOK, retFail := f()
		if retFail != nil {
			if fut.onCatchErr != nil {
				fut.onCatchErr(retFail)
			}
		} else {
			if fut.onThenCB != nil {
				fut.onThenCB(retOK)
			}
		}

		if fut.onComplete != nil {
			fut.onComplete()
		}
	}()

	return fut
}

func (self *future[T]) Then(f func(ret T)) Future[T] {
	self.onThenCB = f
	return self
}

func (self *future[T]) Catch(f func(error)) Future[T] {
	self.onCatchErr = f
	return self
}

func (self *future[T]) Complete(f func()) Future[T] {
	self.onComplete = f
	return self
}
