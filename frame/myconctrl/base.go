package myconctrl

type (
	Golimiter interface {
		Begin()
		End()
		Number() int
	}
)
