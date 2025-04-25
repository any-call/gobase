//go:build darwin

package mylog

import (
	"fmt"
	"os/exec"
)

func OpenLogViewer(logPath string) error {
	cmdStr := fmt.Sprintf(`tell application "Terminal"
    activate
    do script "tail -f '%s'"
end tell`, logPath)

	cmd := exec.Command("osascript", "-e", cmdStr)
	return cmd.Run()
}
