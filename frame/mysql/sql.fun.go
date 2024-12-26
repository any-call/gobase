package mysql

import (
	"fmt"
)

const (
	DateFormatDef = "YYYY-MM-DD HH:MM:SS"
	DateFormatYMD = "%Y-%m-%d"
	DateFormatYM  = "%Y-%m"
)

func Sum[E any](field string, defaultV E, newField string) string {
	return fmt.Sprintf("COALESCE(sum(%s),%v) as %s", field, defaultV, newField)
}

func UNIXSecToDate(field string, format string, newField string) string {
	return fmt.Sprintf("FROM_UNIXTIME(%s, '%s') AS %s", field, format, newField)
}
