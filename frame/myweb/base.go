package myweb

import (
	"net/http"
)

type (
	Binding interface {
		Name() string
		Bind(req *http.Request, obj any) error
	}

	Response interface {
		Write(w http.ResponseWriter, httpcode int, obj any) error
	}
)
