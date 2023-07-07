package mybinddata

// 定义 监听者 标准接口
type Listener interface {
	DataChanged(any)
}

// 定义 Bindable Data 标准接口
type BindData interface {
	AddListener(Listener, any) error
	RemoteListener(Listener) error
}

var ShareBindData BindData = shareBaseBindObj
