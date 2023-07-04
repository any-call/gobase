package mysql

const (
	DateFormatter1 = "%Y-%m-%d %H:%i:%s"
)

//func DateToStr(field string, format string) string {
//	return fmt.Sprintf("date_format(%s,'%s')", field, format)
//}
//
//func StrToDate(dateStr string, format string) string {
//	return fmt.Sprintf("str_to_date('%s','%s')", dateStr, format)
//}

//DATE_FORMAT(convert_tz(check_time,'UTC','Asia/Shanghai'),'%Y-%m-%d %H:%i:%s')  as last_check_time,
//DATE_FORMAT(convert_tz(FROM_UNIXTIME(check_time),'Asia/Shanghai','Asia/Shanghai'),'%Y-%m-%d %H:%i:%s')  as last_check_time1,
//DATE_FORMAT(check_time,'%Y-%m-%d %H:%i:%s')  as last_check_time2,
