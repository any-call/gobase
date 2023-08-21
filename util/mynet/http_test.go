package mynet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func Test_aaaa(t *testing.T) {
	if dnsip, err := GetLocalDNSServer(); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("dns :", dnsip)
	}

	t.Log("ret is :")
}

func BenchmarkGetJson(t *testing.B) {
	var resp map[string]any
	err := GetJson("http://localhost:8081/api/site/sync", nil, 0, func(ret []byte, httpCode int) error {
		if httpCode != http.StatusOK {
			return fmt.Errorf("%d with %s", httpCode, string(ret))
		}

		err := json.Unmarshal(ret, &resp)
		if err != nil {
			return err
		}

		return nil
	}, nil)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ret is :", resp)
}
