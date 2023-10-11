package myweb

import (
	"net/http"
)

func do[REQ, RESP any](r *http.Request, w http.ResponseWriter, req REQ, bindFunc BindFunc[REQ], validateFunc ValidFunc[REQ], handleFunc ServiceFunc[REQ, RESP], writeFunc WriteFunc[RESP]) {
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
