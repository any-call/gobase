package myweb

import "net/http"

var (
	msgpackContentType   = []string{"application/msgpack; charset=utf-8"}
	htmlContentType      = []string{"text/html; charset=utf-8"}
	jsonContentType      = []string{"application/json; charset=utf-8"}
	jsonpContentType     = []string{"application/javascript; charset=utf-8"}
	jsonAsciiContentType = []string{"application/json"}
	protobufContentType  = []string{"application/x-protobuf"}
)

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
