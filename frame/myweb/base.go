package myweb

import (
	"net/http"
)

type (
	//参数绑定
	Binding interface {
		Name() string
		Bind(req *http.Request, obj any) error
	}

	//入参检测
	Valdate interface {
		Check() error
	}

	//响应数据
	Render interface {
		Render(w http.ResponseWriter, httpcode int) error
	}
)
