package mysub

import (
	"testing"
)

func TestBuildSSSubscription(t *testing.T) {
	nodes := []SSNode{
		{
			Name:     "🇭🇰 香港01",
			Server:   "1.2.3.4",
			Port:     31001,
			Method:   "aes-128-gcm",
			Password: "123456",
		},
		{
			Name:     "🇭🇰 香港02",
			Server:   "1.2.3.4",
			Port:     31002,
			Method:   "aes-128-gcm",
			Password: "123456",
		},
	}

	sub := BuildSSSubscription(nodes)

	t.Log("sub is :", sub)
}
