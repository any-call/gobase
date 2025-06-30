package mytime

import (
	"fmt"
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

func BeginningOfDayByLocation(t time.Time, location *time.Location) time.Time {
	year, month, day := t.In(location).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}

// EndOfDay 返回指定日期的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func EndOfDayByLocation(t time.Time, location *time.Location) time.Time {
	year, month, day := t.In(location).Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, location)
}

func CurrBeginningOfDay() time.Time {
	return BeginningOfDay(time.Now())
}

func CurrEndOfDay() time.Time {
	return EndOfDay(time.Now())
}

func UnixNano(nanoSec int64) time.Time {
	return time.Unix(nanoSec/1_000_000_000, nanoSec%1_000_000_000)
}

func HumanizeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d 分钟前", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d 小时前", hours)
	} else if diff < 48*time.Hour {
		return "昨天"
	} else if diff < 72*time.Hour {
		return "前天"
	} else if isSameWeek(t, now) {
		return "本周"
	} else if isLastWeek(t, now) {
		return "上周"
	} else if isSameMonth(t, now) {
		return "本月"
	} else if isSameYear(t, now) {
		months := int(now.Month()) - int(t.Month()) + (now.Year()-t.Year())*12
		return fmt.Sprintf("%d 个月前", months)
	} else {
		years := now.Year() - t.Year()
		return fmt.Sprintf("%d 年前", years)
	}
}

// 判断是否是同一周
func isSameWeek(a, b time.Time) bool {
	yearA, weekA := a.ISOWeek()
	yearB, weekB := b.ISOWeek()
	return yearA == yearB && weekA == weekB
}

// 判断是否是上周
func isLastWeek(a, b time.Time) bool {
	yearA, weekA := a.ISOWeek()
	yearB, weekB := b.ISOWeek()

	if yearA == yearB && weekB-weekA == 1 {
		return true
	}
	if yearB-yearA == 1 && weekA == 52 && weekB == 1 {
		return true
	}
	return false
}

// 判断是否是同一个月
func isSameMonth(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month()
}

// 判断是否是同一年
func isSameYear(a, b time.Time) bool {
	return a.Year() == b.Year()
}
