package mydata

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSafetyOP_OP(t *testing.T) {
	obj := make(map[string]int)
	aa := NewData(obj)
	wg := sync.WaitGroup{}
	wg.Add(2)

	startT := time.Now()
	go func() {
		defer wg.Done()
		for time.Now().Sub(startT).Seconds() < 10 {
			aa.Set(func(m map[string]int) {
				fmt.Println("set ..", time.Now().Sub(startT).Seconds())
				key := time.Now().Format("2006-01-02 15:04:05")
				if _, ok := m[key]; !ok {
					m[key] = 1
				}
			})
		}
	}()

	go func() {
		defer wg.Done()
		for time.Now().Sub(startT).Seconds() < 10 {
			fmt.Println("get len :", len(aa.Get()))
		}
	}()

	wg.Wait()
	for k, v := range obj {
		t.Log("run ok", k, v)
	}

}
