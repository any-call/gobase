package myweb

import (
	"net/http"
)

type (
	Binding interface {
		Name() string
		Bind(req *http.Request, obj any) error
	}
)
