package mysql

import (
	"fmt"
	"github.com/any-call/gobase/util/mycond"
	"strings"
)

type selectBuilder struct {
	from     string
	fields   []string
	joins    string
	whereAnd []string
	whereOr  []string
	group    string
	having   string
	order    string
	limit    int
	offset   int
}

func NewSelectSQL() SelectBuilder {
	return &selectBuilder{}
}

func (self *selectBuilder) Table(query string) SelectBuilder {
	self.from = query
	return self
}

func (self *selectBuilder) Select(column ...string) SelectBuilder {
	if self.fields == nil || len(self.fields) == 0 {
		self.fields = column
	} else {
		self.fields = append(self.fields, column...)
	}

	return self
}

func (self *selectBuilder) Joins(query string, args ...any) SelectBuilder {
	self.joins = prepare(query, args...)
	return self
}

func (self *selectBuilder) Where(query string, args ...any) SelectBuilder {
	currWhere := prepare(query, args...)
	if currWhere != "" {
		if self.whereAnd == nil {
			self.whereAnd = []string{fmt.Sprintf("%s", currWhere)}
		} else {
			self.whereAnd = append(self.whereAnd, fmt.Sprintf("%s", currWhere))
		}
	}

	return self
}

func (self *selectBuilder) Or(query string, args ...any) SelectBuilder {
	currWhere := prepare(query, args...)
	if currWhere != "" {
		if self.whereOr == nil {
			self.whereOr = []string{fmt.Sprintf("%s", currWhere)}
		} else {
			self.whereOr = append(self.whereOr, fmt.Sprintf("%s", currWhere))
		}
	}

	return self
}

func (self *selectBuilder) Group(name string) SelectBuilder {
	self.group = name
	return self
}

func (self *selectBuilder) Having(query string, args ...any) SelectBuilder {
	self.having = prepare(query, args...)
	return self
}

func (self *selectBuilder) Order(v string) SelectBuilder {
	self.order = v
	return self
}

func (self *selectBuilder) PageLimit(page, limit int) SelectBuilder {
	self.limit = limit
	self.offset = self.limit * (page - 1)
	if self.offset < 0 {
		self.offset = 0
	}

	return self
}

func (self *selectBuilder) ToCountSql() string {
	baseSQL := fmt.Sprintf("select count(*) from %s ", self.from)

	if self.joins != "" {
		baseSQL += " " + self.joins
	}

	whereSql := ""
	if self.whereAnd != nil && len(self.whereAnd) > 0 {
		whereSql = strings.Join(self.whereAnd, " and ")
	}

	if self.whereOr != nil && len(self.whereOr) > 0 {
		if whereSql != "" {
			whereSql += " or "
		}

		whereSql += strings.Join(self.whereOr, " or ")
	}

	if whereSql != "" {
		baseSQL += " where " + whereSql
	}

	if self.group != "" {
		baseSQL += " group by  " + self.group
	}

	if self.having != "" {
		baseSQL += " having  " + self.having
	}

	if self.order != "" {
		if self.group != "" {
			baseSQL += " order by  " + self.order
		}
	}

	return baseSQL
}

func (self *selectBuilder) ToSql() string {
	baseSQL := fmt.Sprintf("select %s from %s ", mycond.If(func() bool {
		if self.fields == nil || len(self.fields) == 0 {
			return true
		}
		return false
	}, "*", strings.Join(self.fields, ",")), self.from)

	if self.joins != "" {
		baseSQL += " " + self.joins
	}

	whereSql := ""
	if self.whereAnd != nil && len(self.whereAnd) > 0 {
		whereSql = strings.Join(self.whereAnd, " and ")
	}

	if self.whereOr != nil && len(self.whereOr) > 0 {
		if whereSql != "" {
			whereSql += " or "
		}

		whereSql += strings.Join(self.whereOr, " or ")
	}

	if whereSql != "" {
		baseSQL += " where " + whereSql
	}

	if self.group != "" {
		baseSQL += " group by  " + self.group
	}

	if self.having != "" {
		baseSQL += " having  " + self.having
	}

	if self.order != "" {
		baseSQL += " order by  " + self.order
	}

	if self.limit > 0 && self.offset >= 0 {
		baseSQL = fmt.Sprintf("%s limit %d,%d ", baseSQL, self.offset, self.limit)
	}

	return baseSQL
}
