package mybus

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

const (
	// PublishService - Client service method
	PublishService = "ClientService.PushEvent"
)

// ClientArg - object containing event for client to publish locally
type ClientArg struct {
	Args  []interface{}
	Topic string
}

// Client - object capable of subscribing to a remote event bus
type client struct {
	eventBus EventBus
	address  string
	path     string
	service  *ClientService
}

// NewClient - create a client object with the address and server path
func NewClient(address, path string, bus EventBus) ClientBus {
	if bus == nil {
		bus = NewEventBus()
	}
	clientBus := new(client)
	clientBus.eventBus = bus
	clientBus.address = address
	clientBus.path = path
	clientBus.service = &ClientService{clientBus, &sync.WaitGroup{}, false}
	return clientBus
}

// EventBus - returns the underlying event bus
func (self *client) Bus() EventBus {
	return self.eventBus
}

func (self *client) ClientService() *ClientService {
	return self.service
}

func (self *client) ServerAddr() string {
	return self.address
}

func (self *client) ServerPath() string {
	return self.path
}

func (self *client) doSubscribe(topic string, fn interface{}, serverAddr, serverPath string, subscribeType SubscribeType) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Server not found -", r)
		}
	}()

	rpcClient, err := rpc.DialHTTPPath("tcp", serverAddr, serverPath)
	defer rpcClient.Close()
	if err != nil {
		return fmt.Errorf("dialing: %v", err)
	}
	args := &SubscribeArg{self.address, self.path, PublishService, subscribeType, topic}
	reply := new(bool)
	err = rpcClient.Call(RegisterService, args, reply)
	if err != nil {
		return fmt.Errorf("Register error: %v", err)
	}
	if *reply {
		return self.Bus().Subscribe(topic, fn)
	}

	return fmt.Errorf("rpc call not response")
}

// Subscribe subscribes to a topic in a remote event bus
func (self *client) Subscribe(topic string, fn interface{}, serverAddr, serverPath string) error {
	return self.doSubscribe(topic, fn, serverAddr, serverPath, Subscribe)
}

// SubscribeOnce subscribes once to a topic in a remote event bus
func (self *client) SubscribeOnce(topic string, fn interface{}, serverAddr, serverPath string) error {
	return self.doSubscribe(topic, fn, serverAddr, serverPath, SubscribeOnce)
}

// Start - starts the client service to listen to remote events
func (self *client) Start() error {
	defer func() {
		p := recover()
		if p != nil {
			fmt.Println("start panic", p)
		}
	}()

	var err error
	service := self.service
	if !service.started {
		server := rpc.NewServer()
		if err := server.Register(service); err != nil {
			return err
		}

		server.HandleHTTP(self.path, "/"+self.path)
		l, err := net.Listen("tcp", self.address)
		if err == nil {
			service.wg.Add(1)
			service.started = true
			go http.Serve(l, nil)
		}
	} else {
		err = errors.New("Client service already started")
	}
	return err
}

// Stop - signal for the service to stop serving
func (self *client) Stop() {
	service := self.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}

// ClientService - service object listening to events published in a remote event bus
type ClientService struct {
	client  ClientBus
	wg      *sync.WaitGroup
	started bool
}

// PushEvent - exported service to listening to remote events
func (service *ClientService) PushEvent(arg *ClientArg, reply *bool) error {
	fmt.Println("received event", arg.Topic, arg.Args)
	service.client.Bus().Publish(arg.Topic, arg.Args...)
	*reply = true
	return nil
}
