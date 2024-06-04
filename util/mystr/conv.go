package mystr

import (
	"encoding/hex"
	"strconv"
	"strings"
)

func HexToInt64(hexStr string) (int64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	return strconv.ParseInt(hexStr, 16, 64)
}

func HexToNum[T int | int8 | int16 | int32 | int64](hexStr string) (ret T, err error) {
	var bytes []byte
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr)%2 != 0 { //hex 不能处理奇数长的16进制字符，补0处理
		hexStr = "0" + hexStr
	}
	bytes, err = hex.DecodeString(hexStr)
	if err != nil {
		return
	}
	// 将字节数组转换为 int64
	var intValue T
	for _, b := range bytes {
		intValue = (intValue << 8) + T(b)
	}
	ret = intValue
	return
}
