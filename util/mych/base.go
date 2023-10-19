package mych

type (
	ValidFunc[DATA any]  func(data DATA) bool
	HandleFunc[DATA any] func(data DATA) bool
	Product[DATA any]    interface {
		Producter() chan<- DATA
		Consumer() <-chan DATA
		Send(data DATA)
		SendBy(valid ValidFunc[DATA], data DATA)
		Receive() DATA
		ReceiveBy(handler HandleFunc[DATA])
	}
)
