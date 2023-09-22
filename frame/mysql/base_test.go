package mysql

import "testing"

func TestSelectBuilder_Select(t *testing.T) {
	selectSql := NewSelectSQL()
	tmpSql := selectSql.Table("base_ip_region").Where("start_ip_num > ? and end_ip_num < ? ", 16777216, 19777216).ToSql()
	t.Log("tmpSql :", tmpSql)

	t.Run("prepare", func(t *testing.T) {
		tmpStr := prepare("select * from aa where a between ? and ?", 1, 2, 3, 4)
		t.Log("tmpStr :", tmpStr)
	})
}
