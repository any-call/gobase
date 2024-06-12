package mymail

import "testing"

func TestSend(t *testing.T) {
	err := SendByGmail("luis.giga11@gmail.com", "uhrq cogt ulzi nqun", "156711203@qq.com", "register code", "verify code is :76767", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send ok")
}
