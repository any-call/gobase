package mycmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
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

func ExecbashWithTimeout(
	cmdStr string,
	configFn func(c *exec.Cmd),
	needSudo bool,
	timeout time.Duration,
	args ...string,
) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	list := []string{cmdStr}
	list = append(list, args...)

	command := strings.Join(list, " ")

	var cmd *exec.Cmd
	if needSudo {
		cmd = exec.CommandContext(ctx, "sudo", "bash", "-c", command)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	}

	if configFn != nil {
		configFn(cmd)
	}

	buf, err := cmd.CombinedOutput()
	output = strings.TrimSpace(string(buf))

	if ctx.Err() != nil {
		return output, fmt.Errorf("command timeout: %w", ctx.Err())
	}

	return output, err
}
