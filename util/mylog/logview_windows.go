//go:build windows

package mylog

import (
	"os/exec"
)

func OpenLogViewer(logPath string) error {
	cmd := exec.Command("notepad.exe", logPath)
	return cmd.Start()
}
