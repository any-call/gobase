package mylog

import "testing"

func TestDebug(t *testing.T) {
	SetOptions(WithDisableColor(true))
	Debug("this is test")
	Info("this is info")
	t.Log("ok")
}
