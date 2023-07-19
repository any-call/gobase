package mybus

import (
	"fmt"
	"testing"
)

type MyType int

func (self *MyType) Math(str string) {
	fmt.Println("math :", str, *self)
}

func TestEventBus_Publish(t *testing.T) {
	var aa MyType = 10
	var bb1 MyType = 110
	var bb2 MyType = 120
	var cc MyType = 140
	var dd MyType = 150

	evtBus := NewEventBus()
	if err := evtBus.SubscribeAsync("math", aa.Math); err != nil {
		t.Error(err)
		return
	}

	if err := evtBus.SubscribeAsync("math", bb1.Math); err != nil {
		t.Error(err)
		return
	}

	if err := evtBus.SubscribeAsync("math", bb2.Math); err != nil {
		t.Error(err)
		return
	}

	if err := evtBus.SubscribeOnceAsync("math", cc.Math); err != nil {
		t.Error(err)
		return
	}

	if err := evtBus.SubscribeOnceAsync("math", dd.Math); err != nil {
		t.Error(err)
		return
	}

	evtBus.Publish("math", "this is aa")
	evtBus.WaitAsync()

	t.Log("run ok")
}
