package mynetwork

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	HttpMethodGet  = "GET"
	HttpMethodPost = "POST"
)

const (
	ContentTypeJson           = "application/json"
	ContentTypeFormUrlencoded = "application/x-www-form-urlencoded"
)

type (
	ReqCallback   func(r *http.Request) (isTls bool, timeout time.Duration, err error)
	ParseCallback func(ret []byte, httpCode int) error
)

func DoReq(method string, url string, reqCB ReqCallback, parseCB ParseCallback) (err error) {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	var isTLS bool
	var timeout time.Duration

	if reqCB != nil {
		if isTLS, timeout, err = reqCB(r); err != nil {
			return err
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	if isTLS {
		defaultTransPort := http.DefaultTransport
		defaultTransPort.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = defaultTransPort
	}

	var resp *http.Response = nil
	resp, err = client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if parseCB != nil {
		return parseCB(body, resp.StatusCode)
	}

	return nil
}

func GetJson(url string, inParam any, timeout time.Duration, parseCb ParseCallback) (err error) {
	return DoReq(HttpMethodGet, url, func(r *http.Request) (isTls bool, tout time.Duration, err2 error) {
		r.Header.Add("Content-Type", ContentTypeJson)
		if inParam != nil {
			if b, err := json.Marshal(inParam); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}
		}

		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb)
}

func PostJson(url string, inParam any, timeout time.Duration, parseCb ParseCallback) (err error) {
	return DoReq(HttpMethodPost, url, func(r *http.Request) (isTls bool, tout time.Duration, err2 error) {
		r.Header.Add("Content-Type", ContentTypeJson)
		if inParam != nil {
			if b, err := json.Marshal(inParam); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}
		}

		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb)
}

func GetQuery(url string, param url.Values, timeout time.Duration, parseCb ParseCallback) error {
	newUrl := fmt.Sprintf("%s?%s", url, param.Encode())
	return DoReq(HttpMethodGet, newUrl, func(r *http.Request) (isTls bool, tout time.Duration, err error) {
		r.Header.Add("Content-Type", ContentTypeFormUrlencoded)
		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb)
}

func PostForm(url string, req string, timeout time.Duration, parseCb ParseCallback) error {
	return DoReq(HttpMethodPost, url, func(r *http.Request) (isTls bool, tout time.Duration, err2 error) {
		r.Header.Add("Content-Type", ContentTypeFormUrlencoded)
		if len(req) > 0 {
			b := []byte(req)
			r.Body = io.NopCloser(bytes.NewBuffer(b))
			r.Header.Add("Content-Length", strconv.Itoa(len(b)))
		}

		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb)
}
