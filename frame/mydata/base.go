package mydata

type (
	Data[DATA any] interface {
		Get() DATA
		Set(fn func(DATA))
		TrySet(fn func(DATA)) bool
	}
)
