package mycrypto

import (
	"path/filepath"
	"testing"
)

func TestUnzip(t *testing.T) {
	srcFile := "/Users/luisjin/Documents/whProj/code/bot-api/cmd/adminUI/fyne-cross/dist/darwin-arm64/bot中控管理.zip"
	if err := Unzip(srcFile, filepath.Dir(srcFile)); err != nil {
		t.Error(err)
		return
	}

	t.Log("run ok ")
}
