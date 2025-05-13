//go:build darwin || linux || android

package myos

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func GetDeviceIdentifier() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, "\"")
				if len(parts) > 3 {
					return parts[3], nil
				}
			}
		}
	case "linux":
		cmd := exec.Command("cat", "/etc/machine-id")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil

	case "android":
		cmd := exec.Command("getprop", "ro.serialno")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			// 备选项
			cmd2 := exec.Command("getprop", "ro.boot.serialno")
			output, err = cmd2.Output()
			if err != nil {
				return "", err
			}
		}
		return strings.TrimSpace(string(output)), nil
	}

	return "", fmt.Errorf("unsupported platform:%s", runtime.GOOS)
}
