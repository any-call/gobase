package myxdb

import (
	"io"
	"os"
)

type XDBSearcher struct {
	*Searcher
}

func NewXDBSearcher(filePath string) (ret *XDBSearcher, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	searcher, err := NewWithBuffer(content)
	if err != nil {
		return nil, err
	}

	return &XDBSearcher{searcher}, nil
}

func (self *XDBSearcher) Search(ipaddr string) string {
	str, _ := self.SearchByStr(ipaddr)
	return str
}
