package mysql

import (
	"fmt"
	"strings"
)

type deleteBuilder struct {
	from     string
	whereAnd []string
	whereOr  []string
}

func NewDeleteSQL() DeleteBuilder {
	return &deleteBuilder{}
}

func (self *deleteBuilder) Table(query string) DeleteBuilder {
	self.from = query
	return self
}

func (self *deleteBuilder) Where(query string, args ...any) DeleteBuilder {
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

func (self *deleteBuilder) Or(query string, args ...any) DeleteBuilder {
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

func (self *deleteBuilder) ToSql() string {
	baseSQL := fmt.Sprintf("delete from %s ", self.from)

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

	return baseSQL
}
