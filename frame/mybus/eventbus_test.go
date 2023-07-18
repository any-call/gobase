package mybus

import (
	"fmt"
	"sync"
	"testing"
)

func TestEventBus_Publish(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		fmt.Println("run 1")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("run 2")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("run 3")
	}()

	wg.Wait()
	t.Log("run over")
}
