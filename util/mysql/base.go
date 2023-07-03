package mysql

import "fmt"

const (
	DateFormatter1 = "%Y-%m-%d %H:%i:%s"
)

func DateToStr(dateColumn string, format string) string {
	return fmt.Sprintf("DATE_FORMAT(%s,'%s')", dateColumn, format)
}

//DATE_FORMAT(convert_tz(check_time,'UTC','Asia/Shanghai'),'%Y-%m-%d %H:%i:%s')  as last_check_time,
//DATE_FORMAT(convert_tz(FROM_UNIXTIME(check_time),'Asia/Shanghai','Asia/Shanghai'),'%Y-%m-%d %H:%i:%s')  as last_check_time1,
//DATE_FORMAT(check_time,'%Y-%m-%d %H:%i:%s')  as last_check_time2,
