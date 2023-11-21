package myfuture

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFuture(t *testing.T) {
	Start(func() (string, error) {
		time.Sleep(time.Second * 1)
		return "ok", nil
	}).Then(func(s string) {
		fmt.Println("receive :", s)
		time.Sleep(time.Second * 5)
	}).Complete(func() {
		fmt.Println(" run complete")
	}).Catch(func(err error) {
		fmt.Println(" run err:", err)
	})

	time.Sleep(time.Second * 4)
}
