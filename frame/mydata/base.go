package mydata

type (
	Data[DATA any] interface {
		Get() DATA
		Set(d DATA)
		SetItem(fn func(DATA))
	}
)
