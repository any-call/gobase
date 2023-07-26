package mybus

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

// NetworkBusService - object capable of serving the network bus
type NetworkService struct {
	wg      *sync.WaitGroup
	started bool
}

// NetworkBus - object capable of subscribing to remote event buses in addition to remote event
// busses subscribing to it's local event bus. Compoed of a server and client
type networkBus struct {
	ClientBus
	ServerBus
	service   *NetworkService
	sharedBus EventBus
	address   string
	path      string
}

// NewNetworkBus - returns a new network bus object at the server address and path
func NewNetworkBus(address, path string) NetworkBus {
	bus := new(networkBus)
	bus.sharedBus = NewEventBus()
	bus.ServerBus = NewServer(address, path, bus.sharedBus)
	bus.ClientBus = NewClient(address, path, bus.sharedBus)
	bus.service = &NetworkService{&sync.WaitGroup{}, false}
	bus.address = address
	bus.path = path
	return bus
}

// EventBus - returns wrapped event bus
func (self *networkBus) Bus() EventBus {
	return self.sharedBus
}

func (self *networkBus) ServerAddr() string {
	return self.address
}

func (self *networkBus) ServerPath() string {
	return self.path
}

// Start - helper method to serve a network bus service
func (self *networkBus) Start() error {
	defer func() {
		p := recover()
		if p != nil {
			fmt.Println("start panic", p)
		}
	}()

	var err error
	service := self.service
	clientService := self.ClientBus.ClientService()
	serverService := self.ServerBus.ServerService()
	if !service.started {
		server := rpc.NewServer()
		if err = server.RegisterName("ServerService", serverService); err != nil {
			return err
		}
		if err = server.RegisterName("ClientService", clientService); err != nil {
			return err
		}
		server.HandleHTTP(self.path, "/"+self.path)
		l, e := net.Listen("tcp", self.address)
		if e != nil {
			return fmt.Errorf("listen error: %v", e)
		}
		service.wg.Add(1)
		go http.Serve(l, nil)
	} else {
		err = errors.New("Server bus already started")
	}

	return err
}

// Stop - signal for the service to stop serving
func (self *networkBus) Stop() {
	service := self.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}
