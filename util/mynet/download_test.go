package mynet

import (
	"github.com/any-call/gobase/util/mylog"
	"net/http"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	downloadUrl := "https://dldir1.qq.com/qqfile/qq/PCQQ9.7.17/QQ9.7.17.29225.exe"
	if err := DownloadFile(downloadUrl, "/Users/luisjin/Desktop/xdb/qq.exe", func(r *http.Request) {
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}, func(readLen, totalLen int64) {
		if totalLen > 0 {
			mylog.Debugf("download progress :%.2f", float64(readLen)/float64(totalLen))
		} else {
			mylog.Debug("download size is :", readLen)
		}

	}); err != nil {
		t.Error(err)
		return
	}

	t.Log("download ok ")
}
