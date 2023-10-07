package mysql

import (
	"fmt"
	"github.com/any-call/gobase/util/myconv"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	DateFormatter1 = "%Y-%m-%d %H:%i:%s"
)

type (
	SelectBuilder interface {
		Table(query string) SelectBuilder
		Select(column ...string) SelectBuilder
		Joins(query string, args ...any) SelectBuilder
		Where(query string, args ...any) SelectBuilder
		Or(query string, args ...any) SelectBuilder
		Group(name string) SelectBuilder
		Order(v string) SelectBuilder
		PageLimit(page, limit int) SelectBuilder
		ToCountSql() string
		ToSql() string
	}

	UpdateBuilder interface {
		Table(name string) UpdateBuilder
		Where(query string, args ...any) UpdateBuilder
		Or(query string, args ...any) UpdateBuilder
		Update(column string, v any) UpdateBuilder
		ToSql() string
	}

	InsertBuilder interface {
		Table(name string) InsertBuilder
		Columns(col ...string) InsertBuilder
		AppendValues(v ...any) InsertBuilder
		ToSql() string
	}

	DeleteBuilder interface {
		Table(name string) DeleteBuilder
		Where(query string, args ...any) DeleteBuilder
		Or(query string, args ...any) DeleteBuilder
		ToSql() string
	}
)

// 占位符替换 ？
func prepare(sql string, args ...any) string {
	if args == nil || len(args) == 0 {
		return sql
	}

	rep, _ := regexp.Compile("\\?")
	var index int = 0
	return rep.ReplaceAllStringFunc(sql, func(s string) string {
		if index < len(args) {
			tmpIndex := index
			index++

			directObj := myconv.DirectObj(args[tmpIndex])
			switch reflect.ValueOf(directObj).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return fmt.Sprintf("%d", directObj)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return fmt.Sprintf("%ud", directObj)

			case reflect.Float32, reflect.Float64:
				return fmt.Sprintf("%f", directObj)

			case reflect.String:
				return fmt.Sprintf("'%s'", directObj)

			case reflect.Array, reflect.Slice:
				{
					listArgs := []string{}
					v := reflect.ValueOf(directObj)
					for i := 0; i < v.Len(); i++ {
						if v.Index(i).Kind() == reflect.String {
							listArgs = append(listArgs, fmt.Sprintf("'%v'", v.Index(i)))
						} else {
							listArgs = append(listArgs, fmt.Sprintf("%v", v.Index(i)))
						}
					}

					return "(" + strings.Join(listArgs, ",") + ")"
				}
			case reflect.Struct:
				{
					switch reflect.TypeOf(directObj).String() {
					case "time.Time":
						return fmt.Sprintf("'%s'", directObj.(time.Time).Format("2006-01-02 15:04:05"))
					default:
						return fmt.Sprintf("%v", directObj)
					}
				}

			default:
				return fmt.Sprintf("%v", directObj)
			}
		}

		//fmt.Println("in string", s)
		return s
	})
}
