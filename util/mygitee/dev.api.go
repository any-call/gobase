package mygitee

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
	"net/http"
	"net/url"
	"os"
	"time"
)

type (
	developerApi struct {
		token string
	}
)

func NewDevApi(token string) *developerApi {
	return &developerApi{token: token}
}

func (self *developerApi) ListTags(owner string, repo string, page int) (list []TagInfo, err error) {
	urlStr := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/tags", owner, repo)
	param := url.Values{}
	param.Add("access_token", self.token)
	param.Add("per_page", "100")
	if page > 0 {
		param.Add("page", fmt.Sprintf("%v", page))
	}

	if err = mynet.GetQuery(urlStr, param, time.Second*10, func(ret []byte, httpCode int) error {
		if httpCode != http.StatusOK {
			return errors.New(string(ret))
		}

		if err = json.Unmarshal(ret, &list); err != nil {
			return err
		}
		return nil
	}, nil); err != nil {
		return nil, err
	}

	return list, nil
}

func (self *developerApi) GetZipFile(owner string, repo string, tagName string, zipFilePath string, timeout time.Duration) error {
	urlStr := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/zipball", owner, repo)
	param := url.Values{}
	param.Add("access_token", self.token)
	if len(tagName) > 0 {
		param.Add("ref", tagName)
	}

	return mynet.GetQuery(urlStr, param, timeout, func(ret []byte, httpCode int) error {
		if httpCode != http.StatusOK {
			return errors.New(string(ret))
		}

		return os.WriteFile(zipFilePath, ret, 0644)
	}, nil)
}
