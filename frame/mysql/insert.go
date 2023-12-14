package mysql

import (
	"fmt"
	"gitee.com/any-call/gobase/util/mycond"
	"strings"
)

type insertValue []any

type insertBuilder struct {
	table      string
	columns    []string
	listValues []insertValue
}

func NewInsertSQL() InsertBuilder {
	return &insertBuilder{}
}

func (self *insertBuilder) Table(query string) InsertBuilder {
	self.table = query
	return self
}

func (self *insertBuilder) Columns(col ...string) InsertBuilder {
	self.columns = col
	return self
}

func (self *insertBuilder) AppendValues(v ...any) InsertBuilder {
	if v != nil {
		if self.listValues == nil {
			self.listValues = []insertValue{v}
		} else {
			self.listValues = append(self.listValues, v)
		}
	}

	return self
}

func (self *insertBuilder) ToSql() string {
	baseSQL := fmt.Sprintf("insert into %s(%s) values",
		self.table,
		mycond.If(func() bool {
			if self.columns != nil {
				return true
			}
			return false
		}, strings.Join(self.columns, ","), ""),
	)

	vTemp := make([]string, len(self.columns))
	for i, _ := range vTemp {
		vTemp[i] = "?"
	}

	vStrTemp := fmt.Sprintf("(%s)", strings.Join(vTemp, ","))

	if self.listValues != nil {
		var listValues []string = make([]string, len(self.listValues))
		for j, _ := range self.listValues {
			listValues[j] = prepare(vStrTemp, self.listValues[j]...)
		}

		return baseSQL + strings.Join(listValues, ",")
	}

	return baseSQL
}
