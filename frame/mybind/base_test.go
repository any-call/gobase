package mybind

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
	type myStruct struct {
		Name string
		Sex  int
	}

	var monitorObj map[string]int = make(map[string]int, 0)

	listener := &Dept{
		Name:  "jin",
		Total: 10,
	}

	if err := AddListener(listener, &monitorObj); err != nil {
		t.Error(err)
		return
	}

	go func() {
		var i int
		for {
			i += 5
			SetData(func() {
				monitorObj["name"] = i
			})

			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(time.Second * 10)
	t.Log("run ok ")
}
