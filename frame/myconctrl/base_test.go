package myconctrl

import (
	"fmt"
	"testing"
	"time"
)

func TestGoLimiter_Begin(t *testing.T) {
	limiter := NewGolimiter(10)

	go func() {
		var i int
		for i < 100 {
			limiter.Begin()
			fmt.Println("i = ", i)
			i++
			time.Sleep(time.Millisecond * 10)
			limiter.End()
		}
		t.Log("i over")
	}()

	go func() {
		var m int
		for m < 100 {
			limiter.Begin()
			fmt.Println("m = ", m)
			m++
			time.Sleep(time.Second)
			limiter.End()
		}
		t.Log("m over")
	}()

	time.Sleep(time.Second * 130)
	t.Log("begin over")
}
