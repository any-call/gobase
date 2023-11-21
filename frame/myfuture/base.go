package myfuture

type Future[T any] interface {
	Then(func(T)) Future[T]
	Catch(func(error)) Future[T]
	Complete(func()) Future[T]
}
