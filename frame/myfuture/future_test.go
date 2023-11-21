package myfuture

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFuture(t *testing.T) {
	Start(func() (string, error) {
		time.Sleep(time.Nanosecond)
		panic("this is panic")
		return "OK", nil
	}).Then(func(s string) {
		fmt.Println("receive :", s)
	}).Complete(func() {
		fmt.Println(" run complete")
	}).Catch(func(err error) {
		fmt.Println(" run err:", err)
	})

	//fu1.Cancel()
	//ret, err := fu1.Get()
	//t.Log("ret :", ret, err)
	time.Sleep(time.Second * 2)
}
