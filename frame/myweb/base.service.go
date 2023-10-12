package myweb

import "net/http"

type (
	noReq  struct{}
	noResp struct{}

	NoReqNoRespService         func() error
	ReqNoRespService[REQ any]  func(req REQ) error
	NoReqRespService[RESP any] func() (resp RESP, err error)
)

func noReqNoRespFuncWrap(thenFun NoReqNoRespService) ServiceFunc[noReq, noResp] {
	return func(req noReq) (resp noResp, err error) { err = thenFun(); return }
}

func reqNoRespFuncWrap[REQ any](thenFun ReqNoRespService[REQ]) ServiceFunc[REQ, noResp] {
	return func(req REQ) (resp noResp, err error) { err = thenFun(req); return }
}

func noReqRespFuncWrap[RESP any](thenFun NoReqRespService[RESP]) ServiceFunc[noReq, RESP] {
	return func(req noReq) (resp RESP, err error) { resp, err = thenFun(); return }
}

func do[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, bindFunc BindFunc[REQ], validateFunc ValidFunc[REQ], handleFunc ServiceFunc[REQ, RESP]) {
	if fn := bindFunc; fn != nil {
		if err := fn(r, &req); err != nil {
			WriteBindError(w, err)
			return
		}
	}

	if fn := validateFunc; fn != nil {
		if err := fn(&req); err != nil {
			WriteBindError(w, err)
			return
		}
	}

	if fn := handleFunc; fn != nil {
		if resp, err := handleFunc(req); err != nil {
			WriteServiceError(w, err)
		} else {
			WriteSuccessJSON(w, resp)
		}
	}
}
