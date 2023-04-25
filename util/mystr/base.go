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

func SplitWithRuneLen(s string, length int) []string {
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
