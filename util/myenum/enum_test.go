package myenum

import "testing"

func TestEnum_Name(t *testing.T) {
	callBackMainMenu := NewENum("main_menu", "主菜单")
	t.Log("name:", callBackMainMenu.Name())
	t.Log("value:", callBackMainMenu.Value())
}
