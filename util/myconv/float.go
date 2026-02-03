package myconv

import (
	"fmt"
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

func FloatToPercent[T float32 | float64](v T) string {
	s := fmt.Sprintf("%.6f", v*100) // 最多6位
	s = strings.TrimRight(s, "0")   // 去掉尾部0
	s = strings.TrimRight(s, ".")   // 去掉尾部点
	return s + "%"
}
