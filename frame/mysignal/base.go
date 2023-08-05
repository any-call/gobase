package mysignal

type (
	Signal interface {
		Connect(fn any) (int, error)
		DisConnect(int) error
		Emit(args ...any)
	}
)
