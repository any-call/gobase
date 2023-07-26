package mybus

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mymap"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

// SubscribeType - how the client intends to subscribe
type SubscribeType int

const (
	// Subscribe - subscribe to all events
	Subscribe SubscribeType = iota
	// SubscribeOnce - subscribe to only one event
	SubscribeOnce
)

const (
	// RegisterService - Server subscribe service method
	RegisterService = "ServerService.Register"
)

// SubscribeArg - object to hold subscribe arguments from remote event handlers
type SubscribeArg struct {
	ClientAddr    string
	ClientPath    string
	ServiceMethod string
	SubscribeType SubscribeType
	Topic         string
}

// Server - object capable of being subscribed to by remote handlers
type server struct {
	eventBus    EventBus
	address     string
	path        string
	subscribers *mymap.MultiMap[string, *SubscribeArg]
	service     *ServerService
}

// NewServer - create a new Server at the address and path
func NewServer(address, path string, bus EventBus) ServerBus {
	if bus == nil {
		bus = NewEventBus()
	}

	serverBus := new(server)
	serverBus.eventBus = bus
	serverBus.address = address
	serverBus.path = path
	serverBus.subscribers = mymap.NewMultiMap[string, *SubscribeArg]()
	serverBus.service = &ServerService{serverBus, &sync.WaitGroup{}, false}
	return serverBus
}

// EventBus - returns wrapped event bus
func (self *server) Bus() EventBus {
	return self.eventBus
}

func (self *server) RegisterType(param any) {
	if param != nil {
		gob.Register(param)
	}
}

func (self *server) Service() *ServerService {
	return self.service
}

func (self *server) ServerAddr() string {
	return self.address
}

func (self *server) ServerPath() string {
	return self.path
}

func (self *server) RPCCallback(subscribeArg *SubscribeArg) func(args ...interface{}) {
	return func(args ...interface{}) {
		client, connErr := rpc.DialHTTPPath("tcp", subscribeArg.ClientAddr, subscribeArg.ClientPath)
		defer client.Close()
		if connErr != nil {
			fmt.Errorf("dialing: %v", connErr)
		}
		clientArg := new(ClientArg)
		clientArg.Topic = subscribeArg.Topic
		clientArg.Args = args
		var reply bool
		err := client.Call(subscribeArg.ServiceMethod, clientArg, &reply)
		if err != nil {
			fmt.Println("server dialing client err", clientArg.Args, err)
		}
	}
}

// HasClientSubscribed - True if a client subscribed to this server with the same topic
func (self *server) HasClientSubscribed(arg *SubscribeArg) bool {
	if topicSubscribers, ok := self.subscribers.Values(arg.Topic); ok {
		for _, topicSubscriber := range topicSubscribers {
			if *topicSubscriber == *arg {
				return true
			}
		}
	}
	return false
}

func (self *server) ClientSubscribed() *mymap.MultiMap[string, *SubscribeArg] {
	return self.subscribers
}

// Start - starts a service for remote clients to subscribe to events
func (self *server) Start() error {
	defer func() {
		p := recover()
		if p != nil {
			fmt.Println("start panic", p)
		}
	}()

	var err error
	service := self.service
	if !service.started {
		rpcServer := rpc.NewServer()
		if err = rpcServer.Register(service); err != nil {
			return err
		}

		rpcServer.HandleHTTP(self.path, "/"+self.path)
		l, e := net.Listen("tcp", self.address)
		if e != nil {
			return fmt.Errorf("listen error: %v", e)
		}
		service.started = true
		service.wg.Add(1)
		go http.Serve(l, nil)
	} else {
		err = errors.New("Server bus already started")
	}
	return err
}

// Stop - signal for the service to stop serving
func (self *server) Stop() {
	service := self.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}

// ServerService - service object to listen to remote subscriptions
type ServerService struct {
	serverBus ServerBus
	wg        *sync.WaitGroup
	started   bool
}

// Register - Registers a remote handler to this event bus
// for a remote subscribe - a given client address only needs to subscribe once
// event will be republished in local event bus
func (service *ServerService) Register(arg *SubscribeArg, success *bool) error {
	subscribers := service.serverBus.ClientSubscribed()
	if !service.serverBus.HasClientSubscribed(arg) {
		rpcCallback := service.serverBus.RPCCallback(arg)
		switch arg.SubscribeType {
		case Subscribe:
			if err := service.serverBus.Bus().Subscribe(arg.Topic, rpcCallback); err != nil {
				return err
			}
			break

		case SubscribeOnce:
			if err := service.serverBus.Bus().SubscribeOnce(arg.Topic, rpcCallback); err != nil {
				return err
			}
			break
		}

		subscribers.Insert(arg.Topic, arg)
	}
	*success = true
	return nil
}
