package mysignal

import (
	"fmt"
	"testing"
)

type User struct {
	Name         string
	Sex          bool
	Age          int
	OnAgeChanged Signal
}

func (User *User) AddAge(a int) {
	User.Age = a
	User.OnAgeChanged.Emit(User.Age)
}

func Test_Signal(t *testing.T) {
	user := User{
		Name:         "luis",
		Sex:          false,
		Age:          0,
		OnAgeChanged: NewSignal(),
	}

	fn1 := func(age int) {
		fmt.Println("fn1 is :", age)
	}

	fn2 := func(age int) {
		fmt.Println("fn2 is :", age)
	}

	if _, err := user.OnAgeChanged.Connect(fn1); err != nil {
		t.Error(err)
	}

	if _, err := user.OnAgeChanged.Connect(fn2); err != nil {
		t.Error(err)
	}

	user.AddAge(10)
	user.AddAge(11)

	t.Log("run over")
}
