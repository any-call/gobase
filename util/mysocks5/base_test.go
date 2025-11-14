package mysocks5

import "testing"

func TestConnToSocks5UDP(t *testing.T) {
	_, udpAddr, err := ConnToSocks5UDP(10, "89.222.109.131:12347", func() (userName, password string) {
		return "***", "***"
	}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("udpAddr is :", udpAddr)
}
