package myweb

import (
	"encoding/json"
	"net/http"
)

var (
	htmlContentType = []string{"text/html; charset=utf-8"}
	jsonContentType = []string{"application/json; charset=utf-8"}
)

type kv struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, code, httpCode int, msg string, err error, data any) error {
	dd := kv{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	if err != nil {
		dd.Msg = err.Error()
	}

	body, err := json.Marshal(dd)
	if err != nil {
		return err
	}

	return writeData(w, httpCode, jsonContentType, body)
}

func WriteText(w http.ResponseWriter, code, httpCode int, data string) error {
	return writeData(w, httpCode, htmlContentType, []byte(data))
}

func writeData(w http.ResponseWriter, httpcode int, contentType []string, body []byte) error {
	w.WriteHeader(httpcode)
	writeContentType(w, contentType)
	if bodyAllowedForStatus(httpcode) {
		if _, err := w.Write(body); err != nil {
			return err
		}
	}

	return nil
}

// 设置content type
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
