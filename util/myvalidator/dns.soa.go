package myvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SOA 权威记录起始，
type SOAInfo struct {
	Ns      string //区域主要NS 服务器 域名
	Mbox    string //区域负责人的电子邮件地址或服务器域名
	Serial  uint32 //区域序列号 ：如果附属于此服务器的辅助名称服务器观察到此数字增加，则从服务器将假定该区域已更新并启动
	Refresh uint32 //辅助名称服务器应向主服务器查询SOA记录以检测区域更改的周期秒数。对于小而稳定的区域推荐86400秒（24小时
	Retry   uint32 //如果主服务器没有响应，辅助服务器重新请求SOA记录的秒数，必须小于REFRESH。对于小而稳定的区域推荐7200秒（2小时）
	Expire  uint32 //如果主服务器没有响应，辅助服务器应该停止应答请求的秒数，必须大于REFRESH和RETRY的和。对于小而稳定的区域推荐3600000秒（1000小时）
	Minttl  uint32 //用于计算消极响应缓存的存活时间
}

func NewSOAInfo(srv string) (*SOAInfo, error) {
	list := strings.Fields(srv)
	if len(list) != 7 {
		return nil, errors.New("soa 格式不正确")
	}

	tmpNS := list[0]
	if b := ValidDomain(tmpNS); !b {
		return nil, errors.New("区域NS不是正确的域名")
	}

	tmpMBox := list[1]

	intSerial, err := strconv.Atoi(list[2])
	if err != nil {
		return nil, errors.New("区域序列号是不正确的数字")
	}

	intRefresh, err := strconv.Atoi(list[3])
	if err != nil {
		return nil, errors.New("刷新周期是不正确的秒数")
	}

	intRetry, err := strconv.Atoi(list[4])
	if err != nil {
		return nil, errors.New("重试周期是不正确的秒数")
	}

	//if intRetry >= intRefresh {
	//	return nil, errors.New("辅助服务器重新请求秒数应小于检测周期秒数")
	//}

	intExpire, err := strconv.Atoi(list[5])
	if err != nil {
		return nil, errors.New("timeout是不正确的秒数")
	}

	//if intExpire <= (intRefresh + intRetry) {
	//	return nil, errors.New("辅助服务器重新请求秒数应小于检测周期秒数")
	//}

	intMinimum, err := strconv.Atoi(list[6])
	if err != nil {
		return nil, errors.New("响应缓存时间是不正确的秒数")
	}

	return &SOAInfo{
		Ns:      tmpNS,
		Mbox:    tmpMBox,
		Serial:  uint32(intSerial),
		Refresh: uint32(intRefresh),
		Retry:   uint32(intRetry),
		Expire:  uint32(intExpire),
		Minttl:  uint32(intMinimum),
	}, nil
}

func (self *SOAInfo) Record() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", self.Ns, self.Mbox, self.Serial, self.Refresh, self.Retry, self.Expire, self.Minttl)
}

func ValidSoa(srv string) bool {
	_, err := NewSOAInfo(srv)
	if err != nil {
		return false
	}
	return true
}
