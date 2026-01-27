package myconv

import (
	"strconv"
	"strings"
)

func FloatToString[T float32 | float64](v T, precision int) string {
	if precision < 0 {
		precision = 0
	}

	// 统一转成 float64 处理
	f := float64(v)

	// 不补 0、不用科学计数法
	s := strconv.FormatFloat(f, 'f', -1, 64)

	// 没有小数点，直接返回
	if !strings.Contains(s, ".") {
		return s
	}

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	decPart := parts[1]

	// 精度 = 0，直接返回整数部分
	if precision == 0 {
		return intPart
	}

	// 小数位不超过 precision
	if len(decPart) <= precision {
		return s
	}

	// 超过 precision，截断
	return intPart + "." + decPart[:precision]
}
