package mysub

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type SSNode struct {
	Name     string
	Server   string
	Port     int
	Method   string
	Password string
}

func buildSSURI(n SSNode) string {

	auth := n.Method + ":" + n.Password

	b64 := base64.StdEncoding.EncodeToString([]byte(auth))

	name := url.QueryEscape(n.Name)

	return fmt.Sprintf(
		"ss://%s@%s:%d#%s",
		b64,
		n.Server,
		n.Port,
		name,
	)
}

func BuildSSSubscription(nodes []SSNode) string {

	var lines []string

	for _, n := range nodes {

		lines = append(lines, buildSSURI(n))

	}

	text := strings.Join(lines, "\n")

	return base64.StdEncoding.EncodeToString([]byte(text))
}
