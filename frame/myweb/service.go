package myweb

import (
	"gitee.com/any-call/gobase/util/myvalidator"
	"net/http"
)

func Query(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc))
}

func QueryReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, BindQuery[REQ], validate[REQ], reqNoRespFuncWrap(thenFunc))
}

func QueryResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc))
}

func QueryReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, BindQuery[REQ], validate[REQ], thenFunc)
}

func Body(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc))
}

func BodyReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, BindJson[REQ], validate[REQ], reqNoRespFuncWrap(thenFunc))
}

func BodyResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc))
}

func BodyReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, BindJson[REQ], validate[REQ], thenFunc)
}

func Form(r *http.Request, w http.ResponseWriter, thenFunc NoReqNoRespService) {
	do[noReq, noResp](r, w, noReq{}, nil, nil, noReqNoRespFuncWrap(thenFunc))
}

func FormReq[REQ any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ReqNoRespService[REQ]) {
	do[REQ, noResp](r, w, req, BindForm[REQ], validate[REQ], reqNoRespFuncWrap(thenFunc))
}

func FormResp[RESP any](r *http.Request, w http.ResponseWriter, thenFunc NoReqRespService[RESP]) {
	do[noReq, RESP](r, w, noReq{}, nil, nil, noReqRespFuncWrap(thenFunc))
}

func FormReqResp[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, thenFunc ServiceFunc[REQ, RESP]) {
	do[REQ, RESP](r, w, req, BindForm[REQ], validate[REQ], thenFunc)
}

func validate[REQ any](req *REQ) (err error) {
	if err := myvalidator.Validate(req); err != nil {
		return err
	}

	var chkObj any = *req
	var chkObjPtr any = req

	if v, ok := chkObj.(CheckFun); ok {
		return v.Check()
	}

	if v, ok := chkObjPtr.(CheckFun); ok {
		return v.Check()
	}

	return nil
}
