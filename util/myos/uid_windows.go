//go:build windows

package myos

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func GetDeviceIdentifier() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //hide console window
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		return strings.TrimSpace(lines[1]), nil
	}

	return "", fmt.Errorf("unsupported platform")
}
