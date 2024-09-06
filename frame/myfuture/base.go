package myfuture

type Future interface {
	Then(func()) Future
	Catch(func(error)) Future
	Complete(func()) Future
}
