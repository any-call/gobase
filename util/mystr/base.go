package mystr

import (
	"strings"
)

func Split(src string, flagList []string) []string {
	return strings.FieldsFunc(src, func(r rune) bool {
		res := false
		for _, s := range flagList {
			if string(r) == s {
				res = true
			}
		}
		return res
	})
}

func SplitByLen(s string, length int) []string {
	sLen := len(s)
	if sLen < length {
		return []string{s}
	}

	list := make([]string, 0)
	for i := 0; i < sLen; i = i + length {
		if sLen > (i + length) {
			list = append(list, string(s[i:i+length]))
		} else {
			list = append(list, string(s[i:sLen]))
		}
	}

	return list
}

func SplitRuneByLen(s string, length int) []string {
	rs := []rune(s)
	rsLen := len(rs)
	if rsLen < length {
		return []string{s}
	}

	list := make([]string, 0)
	for i := 0; i < rsLen; i = i + length {
		if rsLen > (i + length) {
			list = append(list, string(rs[i:i+length]))
		} else {
			list = append(list, string(rs[i:rsLen]))
		}
	}

	return list
}

func RemoveSpec(str string, spec string) string {
	return strings.Join(strings.Split(str, spec), "")
}

func RemoveSpace(str string) string {
	return strings.Join(strings.Fields(str), "")
}

// 计算给定src字符串中出现cond的次数
func CalcStrFreq(str, cond string) (n int) {
	src := []rune(str)
	condCount, srcCount := len([]rune(cond)), len(src)
	if condCount > srcCount {
		return -1
	}

	if condCount == 1 {
		for i := 0; i < srcCount; i++ {
			if string(src[i]) == cond {
				n += 1
			}
		}
		return
	}

	for i := 0; i < srcCount; {
		if i+1 >= srcCount || i+condCount > srcCount {
			break
		}

		if string(src[i:i+condCount]) == cond {
			n += 1
			i += condCount
			continue
		}

		i += 1
	}

	return
}

func FormatToByteLength(text string, targetBytes int, leftAlign bool) string {
	currentBytes := len([]byte(text))
	if currentBytes >= targetBytes {
		return text // 如果字符串的字节长度已经大于或等于目标长度，直接返回
	}

	// 计算需要补齐的字节数
	paddingBytes := targetBytes - currentBytes

	// 构建填充字符串（使用空格来填充）
	padding := strings.Repeat(" ", paddingBytes)

	// 根据对齐方式补齐
	if leftAlign {
		return text + padding // 左对齐，在右侧填充
	} else {
		return padding + text // 右对齐，在左侧填充
	}
}

func ByteLength(s string) int {
	return len([]byte(s))
}

// TruncateString 截断字符串，保留前后指定长度，中间用省略字符串代替
func TruncateString(input string, frontLen, backLen int, ellipsis string) string {
	if len(input) <= frontLen+backLen {
		return input
	}

	front := input[:frontLen]
	back := input[len(input)-backLen:]
	return front + ellipsis + back
}
