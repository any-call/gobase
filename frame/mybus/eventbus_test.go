package mybus

import (
	"fmt"
	"testing"
)

func Math(a string, b string) {
	fmt.Println("math :", a, b)
}

func TestEventBus_Publish(t *testing.T) {
	evtBus := NewEventBus()
	if err := evtBus.SubscribeAsync("math", Math); err != nil {
		t.Error(err)
		return
	}

	evtBus.Publish("math", "jin", "luis")

	evtBus.WaitAsync()
	t.Log("run ok")
}
