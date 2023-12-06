package mycmd

import (
	"os/exec"
	"strings"
)

func Exec(cmdStr string, args ...string) (output string, err error) {
	var buf []byte
	buf, err = exec.Command(cmdStr, args...).CombinedOutput()
	output = strings.TrimSpace(string(buf))
	return
}

func Execbash(cmdStr string, args ...string) (output string, err error) {
	list := []string{cmdStr}
	if args != nil {
		list = append(list, args...)
	}

	return Exec("bash", "-c", strings.Join(list, " "))
}
