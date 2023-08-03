package myvalidator

import "time"

// 2006-01-02
func ValidDate(d string) bool {
	_, err := time.Parse("2006-01-02", d)
	if err != nil {
		return false
	}

	return true
}

// 2006-01-02 15:04:05
func ValidDateTime(d string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", d)
	if err != nil {
		return false
	}

	return true
}
