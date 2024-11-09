package mycmd

import (
	"os/exec"
	"strings"
)

func Exec(cmdStr string, configFn func(c *exec.Cmd), args ...string) (output string, err error) {
	var buf []byte
	cmd := exec.Command(cmdStr, args...)
	if configFn != nil {
		configFn(cmd)
	}

	buf, err = cmd.CombinedOutput()
	if buf != nil {
		output = strings.TrimSpace(string(buf))
	}

	return
}

func Execbash(cmdStr string, configFn func(c *exec.Cmd), args ...string) (output string, err error) {
	list := []string{cmdStr}
	if args != nil {
		list = append(list, args...)
	}

	return Exec("bash", configFn, "-c", strings.Join(list, " "))
}
