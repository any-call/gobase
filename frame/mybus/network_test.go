package mybus

import (
	"fmt"
	"testing"
)

type AddParam struct {
	ID   int
	Name string
}

func Test_network(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		serverBus := NewServer(":2020", "/_server_bus_b", nil)
		serverBus.RegisterParamType(AddParam{})
		serverBus.Start()

		fn := func(a []string) {
			fmt.Println("fn1 run ok :", a)
		}

		fn1 := func(a AddParam) {
			fmt.Println("fn1 :", a.ID, a.Name)
		}

		clientBus := NewClient(":2025", "/_client_bus_b", nil)
		clientBus.Start()

		clientBus.Subscribe("topic", fn, ":2020", "/_server_bus_b")
		clientBus.Subscribe("fn1", fn1, ":2020", "/_server_bus_b")
		serverBus.Bus().Publish("topic", []string{"111", "luis.jin"})
		serverBus.Bus().Publish("fn1", AddParam{
			ID:   100,
			Name: "luis.jin",
		})

		clientBus.Stop()
		serverBus.Stop()
	})

	t.Run("client", func(t *testing.T) {
		clientBus := NewClient("localhost:2015", "/_client_bus_", nil)

		eventArgs := make([]interface{}, 1)
		eventArgs[0] = 10

		clientArg := &ClientArg{eventArgs, "topic"}
		reply := new(bool)

		fn := func(a int) {
			if a != 10 {
				t.Fail()
			}
		}

		clientBus.eventBus.Subscribe("topic", fn)
		clientBus.service.PushEvent(clientArg, reply)
		if !(*reply) {
			t.Fail()
		}
	})

	t.Run("newwork", func(t *testing.T) {
		netBusA := NewNetworkBus(":2035", "/net_bus")
		netBusA.Start()

		netBusB := NewNetworkBus(":2036", "/net_bus")
		netBusB.Start()

		fnA := func(a int) {
			fmt.Println("fnA ..", a)
		}
		netBusA.Subscribe("topic", fnA, ":2036", "/net_bus")
		netBusB.Bus().Publish("topic", 2323)
	})
}
