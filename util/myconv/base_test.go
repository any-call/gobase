package myconv

import (
	"testing"
)

type MyStruct struct {
	ID   int
	Name string
	sex  string
}

func (self *MyStruct) MyName() string {
	return "this is a test"
}

func (self MyStruct) MySex() string {
	return "man"
}

func TestStrToNum(t *testing.T) {
	if v, err := StrToNum[int]("123"); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("str to num :", v)
	}

	if v, err := StrToNum[float32]("adfd"); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("str to num :", v)
	}
}

func TestToBool(t *testing.T) {
	if b, err := ToAny[float32]("34.34"); err != nil {
		t.Error(err)
	} else {
		t.Log("b:", b)
	}
}

func Test_StructToMap(t *testing.T) {
	myInfo := MyStruct{
		ID:   23,
		Name: "232",
		sex:  "334343",
	}

	info, err := Struct2Map(myInfo)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("map:", info)
}

func Test_trimFloat(t *testing.T) {
	t.Log(TrimFloat(12.001))
	t.Log(TrimFloat(12.00))
	t.Log(TrimFloat(12.120))
}

func TestFloatToPercent(t *testing.T) {
	t.Log(FloatToPercent(1.30))
	t.Log(FloatToPercent(0.305))
}
