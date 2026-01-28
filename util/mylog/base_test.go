package mylog

import "testing"

func TestDebug(t *testing.T) {
	SetOptions(WithDisableColor(true))
	Debug("this is test")
	Info("this is info")
	DebugKV("enter_conn", "user", "1212", "password", 232, "kk", map[string]any{
		"this is test": 12,
	})
	t.Log("ok")
}
