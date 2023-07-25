package mybus

import (
	"fmt"
	"testing"
)

type MyType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Math(a MyType) {
	fmt.Println("math :", a.ID, a.Name)
}

func TestEventBus_Publish(t *testing.T) {
	evtBus := NewEventBus()
	if err := evtBus.SubscribeAsync("math", Math); err != nil {
		t.Error(err)
		return
	}

	evtBus.Publish("math", MyType{
		ID:   100,
		Name: "luis.jin",
	})

	evtBus.WaitAsync()
	t.Log("run ok")
}
