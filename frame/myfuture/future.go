package myfuture

import (
	"context"
	"fmt"
)

//○ wait ： 同时执行多个，异步
//○ then ：结果回调
//○ catchError : 例外
//○ whenComplete : 执行结束
//○ timeout : 延时执行

type future[T any] struct {
	retOK    T
	retFail  error
	complete bool
	wait     chan bool
	ctx      context.Context
	cancel   func()

	//---
	onThenCB   func(T)
	onCachErr  func(err error)
	onComplete func()
}

func Start[T any](f func() (T, error)) Future[T] {
	if f == nil {
		panic("empty func")
	}

	fut := &future[T]{
		wait: make(chan bool),
	}

	fut.ctx, fut.cancel = context.WithCancel(context.Background())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				//说明函数出了例外
				if fut.onCachErr != nil {
					fut.onCachErr(fmt.Errorf("%v", r))
				}

				if fut.onComplete != nil {
					fut.onComplete()
				}
			}
		}()

		fut.retOK, fut.retFail = f()
		if fut.retFail != nil {
			if fut.onCachErr != nil {
				fut.onCachErr(fut.retFail)
			}
		} else {
			if fut.onThenCB != nil {
				fut.onThenCB(fut.retOK)
			}
		}

		if fut.onComplete != nil {
			fut.onComplete()
		}

		//fut.ret = &result[T]{success, failure}
		//fut.complete = true
		//fut.wait <- true
		//close(fut.wait)
	}()

	return fut
}

func (self *future[T]) Then(f func(ret T)) Future[T] {
	self.onThenCB = f
	return self
}

func (self *future[T]) Catch(f func(error)) Future[T] {
	self.onCachErr = f
	return self
}

func (self *future[T]) Complete(f func()) Future[T] {
	self.onComplete = f
	return self
}

func (self *future[T]) Cancel() {
	self.cancel()
}

func (self *future[T]) Get() (T, error) {
	if self.complete {
		return self.retOK, self.retFail
	}

	select {
	case <-self.wait:
		return self.retOK, self.retFail

	case <-self.ctx.Done():
		self.retFail = self.ctx.Err()
		return self.retOK, self.retFail
	}
}
