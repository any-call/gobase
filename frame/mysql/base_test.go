package mysql

import "testing"

func TestSelectBuilder_Select(t *testing.T) {
	selectSql := NewSelectSQL()
	tmpSql := selectSql.Table("base_ip_region").
		PageLimit(1, 20).
		Order("start_ip desc").
		Where("start_ip_num > ? and end_ip_num < ? ", 16777216, 19777216).
		Where("area_country = ?", "中国").
		Where("area_province = ?", "北京").
		ToSql()

	t.Log("tmpSql :", tmpSql)

}
