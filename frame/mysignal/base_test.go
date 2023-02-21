package mysignal

import (
	"fmt"
	"testing"
)

type User struct {
	Name          string
	Sex           bool
	Age           int
	OnAgeChanged  Signal[int]
	OnNameChanged Signal[string]
}

func (User *User) AddAge(a int) {
	User.Age = a
	User.OnAgeChanged.Emit(User.Age)
}

func (User *User) SetName(a string) {
	User.Name = a
	User.OnNameChanged.Emit(User.Name)
}

func Test_Signal(t *testing.T) {
	user := User{
		Name: "luis",
		Sex:  false,
		Age:  0,
	}

	Connect[int](&(user.OnAgeChanged), func(i int) error {
		fmt.Println("received ageChanged1 :", i)
		return nil
	})

	Connect[int](&(user.OnAgeChanged), func(i int) error {
		fmt.Println("received ageChanged2 :", i)
		return nil
	})

	Connect[string](&(user.OnNameChanged), func(name string) error {
		fmt.Println("received nameChanged :", name)
		return nil
	})

	user.AddAge(10)
	user.AddAge(11)
	user.SetName("luis")
}
