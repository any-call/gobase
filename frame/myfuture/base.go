package myfuture

import "fmt"

type Future[T any] interface {
	Get() Result[T]
	Cancel()
}

type Result[S any] interface {
	Success() S
	Failure() error
}

type result[S any] struct {
	success S
	failure error
}

func (self *result[S]) Success() S {
	return self.success
}

func (self *result[S]) Failure() error {
	return self.failure
}

func (self *result[S]) String() string {
	if self.failure != nil {
		return fmt.Sprintf("%v", self.failure)
	} else {
		return fmt.Sprintf("%v", self.success)
	}
}
