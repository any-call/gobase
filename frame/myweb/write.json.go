package myweb

import (
	"encoding/json"
	"net/http"
)

type writeJson struct {
	data any
}

func NewWriteJson(d any) Render {
	return &writeJson{data: d}
}

func (self *writeJson) Render(w http.ResponseWriter, httpcode int) error {
	body, err := json.Marshal(self.data)
	if err != nil {
		return err
	}

	return writeData(w, httpcode, jsonContentType, body)
}
