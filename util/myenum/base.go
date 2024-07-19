package myenum

type (
	ENum[V any] interface {
		Name() string
		Value() V
	}
)
