package myvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SrvInfo struct {
	Priority int    //优先级
	Weight   int    //权重
	Port     uint16 //端口
	HostName string //主机名
}

func NewSrvInfo(srv string) (*SrvInfo, error) {
	list := strings.Fields(srv)
	if len(list) != 4 {
		return nil, errors.New("srv 格式不正确")
	}

	priority, err := strconv.Atoi(list[0])
	if err != nil {
		return nil, errors.New("srv 不正确的优先级")
	}

	weight, err := strconv.Atoi(list[1])
	if err != nil {
		return nil, errors.New("srv 不正确的权重")
	}

	port, err := strconv.Atoi(list[2])
	if err != nil || port < 0 || port > 65535 {
		return nil, errors.New("srv 不正确的端口")
	}

	if b := ValidDomain(list[3]); !b {
		return nil, errors.New("srv 不正确的主机名")
	}

	return &SrvInfo{
		Priority: priority,
		Weight:   weight,
		Port:     uint16(port),
		HostName: list[3],
	}, nil
}

func (self *SrvInfo) Record() string {
	return fmt.Sprintf("%d %d %d %s", self.Priority, self.Weight, self.Port, self.HostName)
}

func ValidSrv(srv string) bool {
	_, err := NewSrvInfo(srv)
	if err != nil {
		return false
	}
	return true
}
