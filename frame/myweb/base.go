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

	//-------
	PageReq struct {
		Limit int `form:"limit" validate:"min(1,无效的分页数据)"`
		Page  int `form:"page"`
	}

	PageResp[T any] struct {
		Total int64 `json:"total"`
		List  []T   `json:"list"`
	}
)
