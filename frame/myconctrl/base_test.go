package myconctrl

import "testing"

func TestGoLimiter_Begin(t *testing.T) {
	limiter := NewGolimiter(100)

	limiter.Begin()
	t.Log("begin run")
	limiter.End()
	t.Log("end run")
}
