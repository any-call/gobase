package myweb

import (
	"errors"
	"github.com/any-call/gobase/util/mycond"
	"net/http"
)

const defaultMemory = 32 << 20

type formBinding struct {
	tagName string
}

type formPostBinding struct {
	tagName string
}

func NewFormBinding() Binding {
	return formBinding{tagName: "form"}
}

func NewFormBindingEx(tag string) Binding {
	return formBinding{tagName: mycond.If(func() bool {
		if tag == "" {
			return true
		}
		return false
	}, "form", tag)}
}

func (self formBinding) Name() string {
	return "form"
}

func (self formBinding) Bind(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	if err := mapFormByTag(obj, req.Form, self.tagName); err != nil {
		return err
	}
	return nil
}

func NewFormPostBinding() Binding {
	return formPostBinding{tagName: "form"}
}

func NewFormPostBindingEx(tag string) Binding {
	return formPostBinding{tagName: mycond.If(func() bool {
		if tag == "" {
			return true
		}
		return false
	}, "form", tag)}
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (self formPostBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapFormByTag(obj, req.PostForm, self.tagName); err != nil {
		return err
	}
	return nil
}
