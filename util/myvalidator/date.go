package myvalidator

import "time"

// 2006-01-02 15:04:05
func ValidDate(format string, d string) bool {
	_, err := time.Parse(format, d)
	if err != nil {
		return false
	}

	return true
}
