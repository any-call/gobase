package mysignal

type (
	SignalObj interface {
		Connect(fn any) (int, error)
		DisConnect(int) error

		Emit(args ...any)      //同步发射信息
		EmitAsync(args ...any) //异步发射信息
	}
)
