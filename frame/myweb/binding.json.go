package myweb

import (
	"encoding/json"
	"errors"
	"net/http"
)

type jsonBinding struct{}

func NewJsonBinding() Binding {
	return jsonBinding{}
}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(obj)
}
