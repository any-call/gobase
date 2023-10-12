package myweb

import (
	"github.com/any-call/gobase/util/myvalidator"
	"net/http"
)

func Query(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc), nil)
}

func QueryReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, nil, nil, reqNoRespFuncWrap(thenFunc), nil)
}

func QueryResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc), nil)
}

func QueryReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, nil, nil, thenFunc, nil)
}

func Body(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc), nil)
}

func BodyReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, nil, nil, reqNoRespFuncWrap(thenFunc), nil)
}

func BodyResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc), nil)
}

func BodyReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, nil, nil, thenFunc, nil)
}

func Form(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc), nil)
}

func FormReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, nil, nil, reqNoRespFuncWrap(thenFunc), nil)
}

func FormResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc), nil)
}

func FormReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, nil, nil, thenFunc, nil)
}

func validate[REQ any](req *REQ) (err error) {
	if err := myvalidator.Validate(req); err != nil {
		return err
	}

	if v, ok := req.(CheckFun); ok {
		return v.Check()
	}

	if v, ok := (&req).(CheckFun); ok {
		return v.Check()
	}

	return nil
}
