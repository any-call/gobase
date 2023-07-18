package cmd

import (
	"fmt"
	"github.com/any-call/gobase/frame/myevtbus"
	"github.com/any-call/gobase/frame/mynetbus"
	"github.com/any-call/gobase/util/mylist"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/myos"
	"github.com/any-call/gobase/util/myvalidator"
	"testing"
	"time"
)

// 测试并集
func TestList_union(t *testing.T) {
	a := []string{"1", "3", "5"}
	b := []string{"11", "3", "15"}

	c := mylist.Union[string](a, b)
	t.Log("union:", c)
}

func TestList_intersect(t *testing.T) {
	a := []string{"11", "3", "5", "32"}
	b := []string{"11", "32", "15"}

	c := mylist.Intersect[string](a, b)
	t.Log("intersect:", c)
}

func TestList_difference(t *testing.T) {
	a := []string{"11", "3", "5"}
	b := []string{"11", "3", "15"}

	c := mylist.Difference[string](a, b)
	t.Log("union:", c)
}

func Test_mylog(t *testing.T) {
	//opt := mylog.WithFormatter(&mylog.JsonFormatter{IgnoreBasicFields: true})
	//mylog.SetOptions(opt)

	mylog.Debug("this is test")
	mylog.Info("this is test")
	mylog.Warn("this is test")
	mylog.Error("this is test")
	mylog.Panic("this is test")
	mylog.Fatal("this is test")

	t.Log("test ok ")
}

func TestList_removeDuplicItem(t *testing.T) {
	a := []string{"11", "11", "12", "12", "13", "14", "15", "15"}
	a1 := mylist.RemoveDuplicItem[string](a)
	t.Log("a:", a)
	t.Log("a1:", a1)
}

func Test_ValidFqdn(t *testing.T) {
	b1 := myvalidator.ValidFqdn("baidu.com")
	b2 := myvalidator.ValidFqdn("aa.baidu.com")
	t.Log("b1", b1)
	t.Log("b2", b2)
}

func Test_ValidEmail(t *testing.T) {
	b1 := myvalidator.ValidEmail("baidu.com")
	b2 := myvalidator.ValidEmail("12121212@cccc.com")
	t.Log("b1", b1)
	t.Log("b2", b2)
}

func Test_os(t *testing.T) {
	path := "/Users/luisjin/Desktop/ip2Asccode.txt"
	b := myos.IsExistPath(path)
	t.Log("IsExistPath :", b)

	b = myos.IsExistDir(path)
	t.Log("IsExistDir :", b)

	b = myos.IsExistFile(path)
	t.Log("IsExistFile :", b)

	file := myos.Filename("/dfdf/dfdf")
	t.Log("Filename :", file)

	dir := myos.Dir("/dfdf/dfdf")
	t.Log("Dir :", dir)
}

func calculator1(a int, b int) {
	fmt.Printf("%d\n", a+b)
}

func calculator2(a int, b int) {
	fmt.Printf("%d\n", (a+b)*10)
}

func Test_EVTBus(t *testing.T) {
	bus := myevtbus.New()
	if err := bus.Subscribe("calculator", calculator1); err != nil {
		t.Error(err)
		return
	}

	if err := bus.SubscribeAsync("calculator", calculator1, true); err != nil {
		t.Error(err)
		return
	}

	bus.Publish("calculator", 20, 30)
	bus.Unsubscribe("calculator", calculator1)
	bus.Publish("calculator", 50, 30)
	time.Sleep(time.Second * 5)
	t.Log("run ok")
}

func Test_netbus(t *testing.T) {
	serverBus := mynetbus.NewServer(":2020", "/_server_bus_b")
	if err := serverBus.Start(); err != nil {
		t.Error(err)
		return
	}

	clientBus := mynetbus.NewClient(":2025", "/_client_bus_b")
	clientBus.Start()

	clientBus.Subscribe("topic", calculator1, ":2020", "/_server_bus_b")
	clientBus.EventBus().Publish("topic", 20, 30)
	serverBus.EventBus().Publish("topic", 10, 50)

	clientBus.Stop()
	serverBus.Stop()
}

func TestNetworkBus(t *testing.T) {
	networkBusA := mynetbus.NewNetworkBus(":2035", "/_net_bus_A")
	networkBusA.Start()

	networkBusB := mynetbus.NewNetworkBus(":2030", "/_net_bus_B")
	networkBusB.Start()

	fnA := func(a int) {
		if a != 10 {
			t.Fail()
		}
	}
	networkBusA.Subscribe("topic-A", fnA, ":2030", "/_net_bus_B")
	networkBusB.EventBus().Publish("topic-A", 10)

	fnB := func(a int) {
		if a != 20 {
			t.Fail()
		}
	}
	networkBusB.Subscribe("topic-B", fnB, ":2035", "/_net_bus_A")
	networkBusA.EventBus().Publish("topic-B", 20)

	networkBusA.Stop()
	networkBusB.Stop()
}
