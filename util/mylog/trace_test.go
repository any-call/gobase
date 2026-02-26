package mylog

import (
	"testing"
)

func TestDebugf(t *testing.T) {
	trace := NewTrace()
	//trace = trace.With("cid", time.Now().UnixMilli()).With("port", "")
	trace.Debug("enter")
	trace = trace.WithKV("ret", "121")
	trace.Info("end", "ret", "11")
}
