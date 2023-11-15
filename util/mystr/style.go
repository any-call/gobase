package mystr

import (
	"strings"
	"unicode"
)

// 将字符串转成 蛇形风格
func ToSnake(input string) string {
	// Convert.
	result := strings.Map(
		func(r rune) rune {
			if IsWhitespace(r) || '-' == r {
				return '_'
			} else {
				return unicode.ToLower(r)
			}
		}, input)

	// Return
	return result
}

// 将字符串转成 驼峰风格
func ToCamel(input string) string {
	// Convert.
	first := true
	prev := ' '
	result := strings.Map(
		func(r rune) rune {
			if first && (IsWhitespace(prev) || '_' == prev || '-' == prev) {
				first = false
				prev = r
				return unicode.ToLower(r)
			} else if !first && (IsWhitespace(prev) || '_' == prev || '-' == prev) {
				prev = r
				return unicode.ToTitle(r)
			} else if IsWhitespace(r) || '_' == r || '-' == r {
				prev = r
				return -1
			} else {
				prev = r
				return unicode.ToLower(r)
			}
		}, input)

	// Return
	return result
}

func ToProperty(input string) string {
	// Convert.
	result := strings.Map(
		func(r rune) rune {

			if IsWhitespace(r) || '_' == r {
				return '-'
			} else {
				return unicode.ToLower(r)
			}
		}, input)

	// Return
	return result
}

// ToHeaderCase converts the string to 'Header-Case' and returns it.
func ToHeader(input string) string {
	// Convert.
	prev := ' '
	result := strings.Map(
		func(r rune) rune {
			if IsWhitespace(prev) || '_' == prev || '-' == prev {
				prev = r
				return unicode.ToTitle(r)
			} else if IsWhitespace(r) || '_' == r {
				prev = r
				return '-'
			} else {
				prev = r
				return unicode.ToLower(r)
			}
		}, input)

	// Return
	return result
}

func ToPascal(input string) string {
	// Convert.
	prev := ' '
	result := strings.Map(
		func(r rune) rune {
			if IsWhitespace(prev) || '_' == prev || '-' == prev {
				prev = r
				return unicode.ToTitle(r)
			} else if IsWhitespace(r) || '_' == r || '-' == r {
				prev = r
				return -1
			} else {
				prev = r
				return unicode.ToLower(r)

			}
		}, input)

	// Return
	return result
}

// 将字符串转成 Title ：每个字母，首字母大写
func ToTitle(input string) string {
	// Convert.
	prev := ' '
	result := strings.Map(
		func(r rune) rune {

			if IsWhitespace(prev) || '_' == prev || '-' == prev {
				prev = r

				return unicode.ToTitle(r)
			} else {
				prev = r

				return unicode.ToLower(r)

			}
		}, input)

	// Return
	return result
}

// FirstUpper 字符串首字母大写
func ToFirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func ToFirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func IsWhitespace(r rune) bool {
	result := false

	switch r {
	case
		'\u0009', // horizontal tab
		'\u000A', // line feed
		'\u000B', // vertical tab
		'\u000C', // form feed
		'\u000D', // carriage return
		'\u0020', // space
		'\u0085', // next line
		'\u00A0', // no-break space
		'\u1680', // ogham space mark
		'\u180E', // mongolian vowel separator
		'\u2000', // en quad
		'\u2001', // em quad
		'\u2002', // en space
		'\u2003', // em space
		'\u2004', // three-per-em space
		'\u2005', // four-per-em space
		'\u2006', // six-per-em space
		'\u2007', // figure space
		'\u2008', // punctuation space
		'\u2009', // thin space
		'\u200A', // hair space
		'\u2028', // line separator
		'\u2029', // paragraph separator
		'\u202F', // narrow no-break space
		'\u205F', // medium mathematical space
		'\u3000': // ideographic space
		result = true
	default:
		result = false
	}

	return result
}
