package myweb

import (
	"errors"
	"github.com/any-call/gobase/util/mycond"
	"net/http"
)

type queryBinding struct{ tagName string }

func NewQueryBinding() Binding {
	return queryBinding{tagName: "form"}
}

func NewQueryBindingEx(tag string) Binding {
	return queryBinding{tagName: mycond.If(func() bool {
		if tag == "" {
			return true
		}
		return false
	}, "form", tag)}
}

func (queryBinding) Name() string {
	return "query"
}

func (self queryBinding) Bind(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	values := req.URL.Query()
	if err := mapFormByTag(obj, values, self.tagName); err != nil {
		return err
	}

	return nil
}
