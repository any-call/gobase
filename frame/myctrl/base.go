package myctrl

type (
	Golimiter interface {
		Begin()
		End()
		Number() int
	}
)
