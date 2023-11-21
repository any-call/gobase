package myfuture

type Future[T any] interface {
	Get() (T, error)
	Cancel()
	Then(func(T)) Future[T]
	Catch(func(error)) Future[T]
	Complete(func()) Future[T]
}

//type Result[S any] interface {
//	Success() S
//	Failure() error
//	String() string
//}
//
//type result[S any] struct {
//	success S
//	failure error
//}
//
//func (self *result[S]) Success() S {
//	return self.success
//}
//
//func (self *result[S]) Failure() error {
//	return self.failure
//}
//
//func (self *result[S]) String() string {
//	if self.failure != nil {
//		return fmt.Sprintf("%v", self.failure)
//	} else {
//		return fmt.Sprintf("%v", self.success)
//	}
//}
