package mysql

import (
	"testing"
	"time"
)

func TestSelectBuilder_Select(t *testing.T) {
	selectSql := NewSelectSQL()
	selectSql.Table("base_ip_region").
		PageLimit(1, 20).
		Order("start_ip desc").
		Where("start_ip_num > ? and end_ip_num < ? ", 16777216, 19777216).
		Where("area_country = ?", "中国").
		Or("area_province = ?", "北京").
		Where("check_time > ?", time.Now()).
		Group("area_continent").
		Select("area_continent", "count(*) as aa")

	t.Log("tmpSql select:", selectSql.ToSql())
	t.Log("tmpSql count:", selectSql.ToCountSql())

	selectSql1 := NewSelectSQL()
	tmpsql1 := selectSql1.Table("manager_user as a").Joins("left join system_role b on a.role_id = b.id").
		Select("a.id as user_id", "a.user_name", "a.role_id", "b.name as role_name").
		Where("b.name in ?", []string{"超级管理员", "测试"}).
		ToSql()
	t.Log("tmpSql1 :", tmpsql1)

}

func TestUpdateBuilder(t *testing.T) {
	updateSql := NewUpdateSQL()
	tmpSql := updateSql.Table("system_role").
		Where("description = ? ", "esse11").
		Update("name", "test").Update("status", 1).
		ToSql()

	t.Log("tmpSql :", tmpSql)

}

func TestInsertBuilder(t *testing.T) {
	insertSql := NewInsertSQL()
	tmpSql := insertSql.Table("task_manager_alert").Columns("task_id", "type_id", "uuid").
		AppendValues(1, 1, "test1").AppendValues(2, 3, "test", "app").
		ToSql()

	t.Log("tmpSql :", tmpSql)

}

func TestDeleteBuilder(t *testing.T) {
	deleteSQL := NewDeleteSQL()
	tmpSql := deleteSQL.Table("black_remoteip_info").Where("id in ?", []int{1, 2, 3}).
		ToSql()

	t.Log("tmpSql :", tmpSql)
}
