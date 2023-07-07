package mybinddata

import (
	"fmt"
	"testing"
	"time"
)

type Dept struct {
	Name  string
	Total int
}

func (d Dept) DataChanged(val any) {
	fmt.Println("data change to :", val)

}

func Test_bind(t *testing.T) {
	var monitorObj map[string]int = make(map[string]int, 0)
	listener := &Dept{
		Name:  "jin",
		Total: 10,
	}

	if err := ShareBindData.AddListener(listener, &monitorObj); err != nil {
		t.Error(err)
		return
	}

	go func() {
		var i int
		for {
			i += 5
			monitorObj["name"] = i
			//fmt.Println("set monitorObj :", monitorObj)
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 10)
	t.Log("run ok ")
}
