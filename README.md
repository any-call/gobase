# gobase  基于go ,封装一些常的框架 及功能库， 会持续保持更新

## util  功能库
## frame 基础框架

### 部分已封装库介绍
### mylist : 用于封装 切片的常用功能
### mysignal :基于QT 的信号 与 槽 的思路来封装，对象间通信号，可以定义多个信号,对象间的通信则通过信号与槽建立联接来实现通信。

#### Demo 如下:

##### type User struct {
##### Name          string
##### Sex           bool
##### Age           int
##### OnAgeChanged  Signal[int]
##### OnNameChanged Signal[string]
##### }

##### func (User *User) AddAge(a int) {
#####   User.Age = a
#####   User.OnAgeChanged.Emit(User.Age)
#####  }

##### func (User *User) SetName(a string) {
##### User.Name = a
##### User.OnNameChanged.Emit(User.Name)
##### }

##### func Test_Signal(t *testing.T) {
##### user := User{
##### Name: "luis",
##### Sex:  false,
##### Age:  0,
##### }

#####	Connect[int](&(user.OnAgeChanged), func(i int) error {
#####		fmt.Println("received ageChanged1 :", i)
#####		return nil
#####	})

#####	Connect[int](&(user.OnAgeChanged), func(i int) error {
#####		fmt.Println("received ageChanged2 :", i)
#####		return nil
#####	})

#####	Connect[string](&(user.OnNameChanged), func(name string) error {
#####		fmt.Println("received nameChanged :", name)
#####		return nil
#####	})

#####	user.AddAge(10)
#####	user.AddAge(11)
#####	user.SetName("luis")
##### }
