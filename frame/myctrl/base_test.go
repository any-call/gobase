package myctrl

import (
	"fmt"
	"testing"
	"time"
)

func TestGoLimiter_Begin(t *testing.T) {
	limiter := NewGolimiter(1)

	go func() {
		var i int
		for i < 100 {
			limiter.Begin()
			fmt.Println("i = ", i)
			i++
			time.Sleep(time.Millisecond * 100)
		}
		t.Log("i over")
	}()

	go func() {
		var m int
		for m < 100 {
			fmt.Println("m = ", m)
			m++
			time.Sleep(time.Millisecond * 10)
			limiter.End()
		}
		t.Log("m over")
	}()

	time.Sleep(time.Second * 130)
	t.Log("begin over")
}

func Test_Delay(t *testing.T) {
	fn := func(n int) {
		fmt.Println("enter : ", n)
	}

	var m int = 100
	DelayExec(time.Second*5, func() {
		fn(m)
	})

	time.Sleep(time.Second * 10)
}
