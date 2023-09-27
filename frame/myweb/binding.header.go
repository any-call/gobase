package myweb

import (
	"github.com/any-call/gobase/util/mycond"
	"net/http"
	"net/textproto"
	"reflect"
)

type headerBinding struct{ tagName string }

func NewHeaderBind() Binding {
	return headerBinding{tagName: "header"}
}

func NewHeaderBindEx(tag string) Binding {
	return headerBinding{tagName: mycond.If(func() bool {
		if tag == "" {
			return true
		}
		return false
	}, "header", tag)}
}

func (headerBinding) Name() string {
	return "header"
}

func (self headerBinding) Bind(req *http.Request, obj any) error {

	if err := mappingByPtr(obj, headerSource(req.Header), self.tagName); err != nil {
		return err
	}

	return nil
}

type headerSource map[string][]string

var _ setter = headerSource(nil)

func (hs headerSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (bool, error) {
	return setByForm(value, field, hs, textproto.CanonicalMIMEHeaderKey(tagValue), opt)
}
