package mysql

import (
	"fmt"
	"strings"
)

type updateBuilder struct {
	from         string
	updateFileds []string
	whereAnd     []string
	whereOr      []string
}

func NewUpdateSQL() UpdateBuilder {
	return &updateBuilder{}
}

func (self *updateBuilder) Table(query string) UpdateBuilder {
	self.from = query
	return self
}

func (self *updateBuilder) Where(query string, args ...any) UpdateBuilder {
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

func (self *updateBuilder) Or(query string, args ...any) UpdateBuilder {
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

func (self *updateBuilder) Update(column string, v any) UpdateBuilder {
	if column != "" {
		currUpdate := prepare(fmt.Sprintf("%s = ?", column), v)
		if self.updateFileds == nil {
			self.updateFileds = []string{currUpdate}
		} else {
			self.updateFileds = append(self.updateFileds, currUpdate)
		}
	}

	return self
}

func (self *updateBuilder) ToSql() string {
	baseSQL := fmt.Sprintf("update %s set ", self.from)

	if self.updateFileds != nil {
		baseSQL += " " + strings.Join(self.updateFileds, ",")
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

	return baseSQL
}
