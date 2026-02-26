package mylog

import (
	"testing"
	"time"
)

func TestDebugf(t *testing.T) {
	trace := NewTrace()
	//trace = trace.With("cid", time.Now().UnixMilli()).With("port", "")
	trace.Debug("enter")
	trace = trace.WithKV("ret", "121")
	trace.Info("end", "ret", "11")

	t.Log(NewCID())
	time.Sleep(time.Second)
	t.Log(NewCID())
}
