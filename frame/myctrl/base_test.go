package myctrl

import (
	"fmt"
	"testing"
	"time"
)

func TestGoLimiter_Begin(t *testing.T) {
	limiter := NewGolimiter(1)

	fmt.Println("test objFun is :", ObjFun(func() string {
		return "123"
	}))
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

func Test_TimeExec(t *testing.T) {
	fn := func(n int) {
		fmt.Println("enter : ", n, time.Now().Second())
		time.Sleep(time.Second * 2)
	}

	var m int = 100
	go TimerExec(time.Second, func() {
		go fn(m)
	})

	time.Sleep(time.Minute)
}
