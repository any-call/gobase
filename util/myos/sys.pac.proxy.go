package myos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type (
	PacProxyConfig struct {
		Enable bool   `json:"enable"` // 是否启用 PAC 代理
		PacURL string `json:"pacURL"` // PAC 代理地址
	}
)

func (self PacProxyConfig) ToString() string {
	ret, _ := json.Marshal(self)
	return string(ret)
}

func SetPacProxy(pacUrl string) error {
	if runtime.GOOS == "windows" {
		cmds := []*exec.Cmd{
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "AutoConfigURL", "/t", "REG_SZ", "/d", pacUrl, "/f"),
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	} else if runtime.GOOS == "darwin" {
		cmds := []*exec.Cmd{
			exec.Command("networksetup", "-setautoproxyurl", "Wi-Fi", pacUrl),
			exec.Command("networksetup", "-setautoproxystate", "Wi-Fi", "on"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}

func ResetPacProxy() error {
	if runtime.GOOS == "windows" {
		cmds := []*exec.Cmd{
			exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "AutoConfigURL", "/f"),
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	} else if runtime.GOOS == "darwin" {
		cmds := []*exec.Cmd{
			exec.Command("networksetup", "-setautoproxystate", "Wi-Fi", "off"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}

func GetPacProxy() (*PacProxyConfig, error) {
	if runtime.GOOS == "windows" {
		// 获取 PAC 代理 URL
		cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "AutoConfigURL")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(output), "\n")
		pacURL := ""
		for _, line := range lines {
			if strings.Contains(line, "AutoConfigURL") {
				parts := strings.Fields(line)
				if len(parts) > 2 {
					pacURL = parts[len(parts)-1]
				}
			}
		}

		// 检查是否启用了代理
		cmd = exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable")
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
		enabled := strings.Contains(string(output), "0x1")

		return &PacProxyConfig{
			Enable: enabled,
			PacURL: pacURL,
		}, nil
	} else if runtime.GOOS == "darwin" {
		// macOS 获取 PAC 代理
		cmd := exec.Command("networksetup", "-getautoproxyurl", "Wi-Fi")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		ret := &PacProxyConfig{}
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "Enabled":
					if value == "No" {
						ret.Enable = false
					} else {
						ret.Enable = true
					}
					break
				case "URL":
					ret.PacURL = value
					break
				}
			}
		}

		return ret, nil
	}

	return nil, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}
