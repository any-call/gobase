package mybus

import (
	"fmt"
	"testing"
)

func Test_network(t *testing.T) {
	type MyType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	t.Run("server", func(t *testing.T) {
		serverBus := NewServer(":2020", "/_server_bus_b", nil)
		serverBus.Start()

		fn := func(a MyType) {
			fmt.Println("fn1 run ok :", a.ID, a.Name)
		}

		clientBus := NewClient(":2025", "/_client_bus_b", nil)
		clientBus.Start()

		clientBus.Subscribe("topic", fn, ":2020", "/_server_bus_b")
		//clientBus.Bus().Publish("topic", MyType{
		//	ID:   50,
		//	Name: "luis.jin",
		//})

		//serverBus.Bus().Publish("topic", MyType{
		//	ID:   2,
		//	Name: "king.king",
		//})

		clientBus.Stop()
		serverBus.Stop()
	})

	//t.Run("server", func(t *testing.T) {
	//
	//})
	//
	//t.Run("network", func(t *testing.T) {
	//
	//})
}
