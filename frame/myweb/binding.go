package myweb

import (
	"encoding/json"
	"errors"
	"net/http"
)

const defaultMemory = 32 << 20

func BindJson(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(obj)
}

func BindQuery(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	values := req.URL.Query()
	if err := mapFormByTag(obj, values, "form"); err != nil {
		return err
	}

	return nil
}

func BindForm(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	if err := mapFormByTag(obj, req.Form, "form"); err != nil {
		return err
	}
	return nil
}

func BindPostForm(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapFormByTag(obj, req.PostForm, "form"); err != nil {
		return err
	}
	return nil
}

func BindHeader(req *http.Request, obj any) error {
	if req == nil {
		return errors.New("invalid request")
	}

	if err := mappingByPtr(obj, headerSource(req.Header), "header"); err != nil {
		return err
	}

	return nil
}
