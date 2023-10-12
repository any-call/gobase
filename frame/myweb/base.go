package myweb

import (
	"net/http"
)

type (
	//参数绑定
	BindFunc[REQ any] func(ctx *http.Request, req *REQ) (err error)

	//入参检测
	ValidFunc[REQ any] func(req *REQ) (err error)

	CheckFun interface {
		Check() error
	}

	//业务处理
	ServiceFunc[REQ, RESP any] func(req REQ) (resp RESP, err error)

	//写包
	WriteFunc[RESP any] func(w http.ResponseWriter, httpCode int, resp RESP, err error)
)
