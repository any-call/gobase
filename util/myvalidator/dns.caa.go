package myvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CAA，全称Certificate Authority Authorization，即证书颁发机构授权
type CAAInfo struct {
	Flags    uint8  //认证机构限制标志
	Tag      string //证书属性标签
	HostName string //证书颁发机构、策略违规报告邮件地址等
}

func NewCAAInfo(caa string) (*CAAInfo, error) {
	list := strings.Fields(caa)
	if len(list) != 3 {
		return nil, errors.New("CAA 格式不正确")
	}

	flags, err := strconv.Atoi(list[0])
	if err != nil || flags < 0 || flags > 255 {
		return nil, errors.New("CAA 不正确的标识值")
	}

	if list[1] != "issue" &&
		list[1] != "issuewild" &&
		list[1] != "iodef" {
		return nil, errors.New("CAA 不正确的标签")
	}

	return &CAAInfo{
		Flags:    uint8(flags),
		Tag:      list[1],
		HostName: list[2],
	}, nil
}

func (self *CAAInfo) Record() string {
	return fmt.Sprintf("%d %s %s", self.Flags, self.Tag, self.HostName)
}

func ValidCAA(caa string) bool {
	_, err := NewCAAInfo(caa)
	if err != nil {
		return false
	}
	return true
}
