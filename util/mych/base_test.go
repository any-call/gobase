package mych

import (
	"testing"
	"time"
)

func TestNewProduct(t *testing.T) {
	productTimer := NewProduct[*int](100)

	go func() {
		for i := 0; i < 100; i++ {
			productTimer.Send(&i)
			time.Sleep(time.Second)
		}
	}()

	go productTimer.ReceiveBy(func(data *int) bool {
		t.Log("received time:", *data)
		return true
	})

	time.Sleep(time.Minute * 5)
	t.Log("run ok")
}
