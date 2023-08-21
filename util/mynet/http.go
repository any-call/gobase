package mynet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

func NewLongClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			Proxy:             http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // tcp连接超时时间
				KeepAlive: 60 * time.Second, // 保持长连接的时间
			}).DialContext, // 设置连接的参数
			MaxIdleConns:          100,              // 最大空闲连接
			MaxConnsPerHost:       100,              //每个host保持的空闲连接数
			MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
			ExpectContinueTimeout: 30 * time.Second, // 等待服务第一响应的超时时间
			IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
		},
	}
}

func DoReq(method string, url string, reqCB ReqCallback, parseCB ParseCallback, client *http.Client) (err error) {
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

	if client == nil {
		client = &http.Client{
			Timeout: timeout,
		}

		if isTLS {
			defaultTransPort := http.DefaultTransport
			defaultTransPort.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			client.Transport = defaultTransPort
		}
	} else {
		client.Timeout = timeout
		if isTLS {
			if _, ok := client.Transport.(*http.Transport); ok {
				client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
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

func GetJson(url string, inParam any, timeout time.Duration, parseCb ParseCallback, client *http.Client) (err error) {
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
	}, parseCb, client)
}

func PostJson(url string, inParam any, timeout time.Duration, parseCb ParseCallback, client *http.Client) (err error) {
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
	}, parseCb, client)
}

func GetQuery(url string, param url.Values, timeout time.Duration, parseCb ParseCallback, client *http.Client) error {
	newUrl := fmt.Sprintf("%s?%s", url, param.Encode())
	return DoReq(HttpMethodGet, newUrl, func(r *http.Request) (isTls bool, tout time.Duration, err error) {
		r.Header.Add("Content-Type", ContentTypeFormUrlencoded)
		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb, client)
}

func PostForm(url string, values url.Values, timeout time.Duration, parseCb ParseCallback, client *http.Client) error {
	return DoReq(HttpMethodPost, url, func(r *http.Request) (isTls bool, tout time.Duration, err2 error) {
		r.Header.Add("Content-Type", ContentTypeFormUrlencoded)
		if values != nil {
			tmpStr := values.Encode()
			r.Body = io.NopCloser(strings.NewReader(tmpStr))
			r.Header.Add("Content-Length", strconv.Itoa(len(tmpStr)))
		}

		tout = timeout
		if strings.HasPrefix(strings.ToLower(url), "https") {
			isTls = true
		}
		return isTls, tout, nil
	}, parseCb, client)
}
