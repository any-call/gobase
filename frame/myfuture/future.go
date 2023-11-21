package myfuture

import (
	"context"
	"fmt"
)

type future[T any] struct {
	ret      *result[T]
	complete bool
	wait     chan bool
	ctx      context.Context
	cancel   func()
}

func NewFuture[T any](f func() (T, error)) Future[T] {
	fut := &future[T]{
		wait: make(chan bool),
	}
	fut.ctx, fut.cancel = context.WithCancel(context.Background())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("future recovered: ", r)
			}
		}()
		success, failure := f()
		fut.ret = &result[T]{success, failure}
		fut.complete = true
		fut.wait <- true
		close(fut.wait)
	}()

	return fut
}

func (self *future[T]) Cancel() {
	self.cancel()
}

func (self *future[T]) Get() Result[T] {
	if self.complete {
		return self.ret
	}

	select {
	case <-self.wait:
		return self.ret

	case <-self.ctx.Done():
		return &result[T]{failure: self.ctx.Err()}
	}
}
