package mytime

import (
	"time"
)

const (
	FormatYYYYMMDDHHMISS = "2006-01-02 15:04:05"
)

var (
	CstZone     = time.FixedZone("CST", 8*3600)
	CstSh, _    = time.LoadLocation("Asia/Shanghai") //上海
	CstLocal, _ = time.LoadLocation("Local")
	CstUTC, _   = time.LoadLocation("UTC")

	TimeTemplate1 = "2006-01-02 15:04:05"
	TimeTemplate2 = "2006/01/02 15:04:05"
	TimeTemplate3 = "2006-01-02"
	TimeTemplate5 = "20060102 150405"
	TimeTemplate6 = "20060102"
)

// 给定年月值,获取该月有多少天
func GetDaysInMonth(yearInt, month int) (days int) {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		days = 31
		return
	case 4, 6, 9, 11:
		days = 30
		return
	case 2:
		if (yearInt%4 == 0 && yearInt%100 != 0) || yearInt%400 == 0 {
			days = 29
			return
		}
		days = 28
		return
	}

	return 31
}

func TruncateMillSec(in time.Time) time.Time {
	return in.Truncate(time.Second)
}

func TruncateSec(in time.Time) time.Time {
	return in.Truncate(time.Minute)
}

func TruncateMinute(in time.Time) time.Time {
	return in.Truncate(time.Hour)
}

func BeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回指定日期的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}
