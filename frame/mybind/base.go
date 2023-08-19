package mybind

// 定义 监听者 标准接口
type Listener interface {
	DataChanged(any)
}

// 定义 Bindable Data 标准接口
type BindData interface {
	AddListener(Listener, any) error
	RemoteListener(Listener) error
	SetState(fn func())
}

func AddListener(l Listener, data any) error {
	return shareBaseBindObj.AddListener(l, data)
}

func RemoteListener(l Listener) error {
	return shareBaseBindObj.RemoteListener(l)
}

func SetState(fn func()) {
	shareBaseBindObj.SetState(fn)
}
