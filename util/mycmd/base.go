package mycmd

import (
	"os/exec"
	"strings"
)

func Exec(cmdStr string, configFn func(c *exec.Cmd), needSudo bool, args ...string) (output string, err error) {
	var buf []byte
	if needSudo {
		// 构建 sudo + 原命令 + 参数
		fullArgs := append([]string{cmdStr}, args...)
		cmd := exec.Command("sudo", fullArgs...)
		if configFn != nil {
			configFn(cmd)
		}
		buf, err = cmd.CombinedOutput()
	} else {
		cmd := exec.Command(cmdStr, args...)
		if configFn != nil {
			configFn(cmd)
		}
		buf, err = cmd.CombinedOutput()
	}

	if buf != nil {
		output = strings.TrimSpace(string(buf))
	}

	return
}

func Execbash(cmdStr string, configFn func(c *exec.Cmd), needSudo bool, args ...string) (output string, err error) {
	list := []string{cmdStr}
	if args != nil {
		list = append(list, args...)
	}

	return Exec("bash", configFn, needSudo, "-c", strings.Join(list, " "))
}
