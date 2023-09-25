package mysql

import (
	"fmt"
	"github.com/any-call/gobase/util/myconv"
	"reflect"
	"regexp"
)

type (
	SelectBuilder interface {
		Table(query string) SelectBuilder
		Select(column ...string) SelectBuilder
		Joins(query string, args ...any) SelectBuilder
		InnerJoins(query string, args ...any) SelectBuilder
		Where(query string, args ...any) SelectBuilder
		Group(name string) SelectBuilder
		Order(v string) SelectBuilder
		PageLimit(page, limit int) SelectBuilder
		ToSql() string
	}

	UpdateBuilder interface {
		Table(name string) UpdateBuilder
		Where(query string, args ...any) UpdateBuilder
		Update(column string, v any) UpdateBuilder
		ToSql() string
	}

	InsertBuilder interface {
		Table(name string) InsertBuilder
		Where(query string, args ...any) InsertBuilder
		Columns(col ...string) InsertBuilder
		Values(v ...[]any) InsertBuilder
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
			default:
				return fmt.Sprintf("%v", directObj)
			}
		}

		//fmt.Println("in string", s)
		return s
	})
}
