package mymail

import "testing"

func TestSendEmail(t *testing.T) {
	t.Log(SendEmail("luis.giga11@gmail.com", "156711203@qq.com", "test mail", "this is a email", 0))
}
