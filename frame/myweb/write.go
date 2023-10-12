package myweb

import (
	"net/http"
)

var (
	HtmlContentType = []string{"text/html; charset=utf-8"}
	JsonContentType = []string{"application/json; charset=utf-8"}
)

type kv struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func WriteJson[RESP any](w http.ResponseWriter, httpCode int, resp RESP, err error) {
	if httpCode != http.StatusOK {
		writeJSON(w, httpCode, httpCode, "", err, resp)
	} else {
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, http.StatusOK, "", err, nil)
		} else {
			writeJSON(w, 0, http.StatusOK, "", nil, resp)
		}
	}
}

func WriteBindError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, http.StatusBadRequest, "", err, nil)
}

func WriteServiceError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, http.StatusInternalServerError, "", err, nil)
}

func WriteSuccessJSON(w http.ResponseWriter, data any) {
	writeJSON(w, 0, http.StatusOK, "", nil, data)
}
