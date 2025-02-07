package myos

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type (
	SocksProxyConfig struct {
		Server string `json:"server"` // 服务器地址 (如 127.0.0.1)
		Port   int    `json:"port"`   // 端口号 (如 1080)
		Enable bool   `json:"enable"`
		//Auth     bool   // 是否需要身份验证 --目录命令设置不了
		//Username string // 用户名 (如果需要) --目录命令设置不了
		//Password string // 密码 (如果需要) --目录命令设置不了
	}
)

func (self SocksProxyConfig) ToString() string {
	ret, _ := json.Marshal(self)
	return string(ret)
}

func SetSocksProxy(server string, port int) error {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// macOS 设置 HTTP 代理
		cmds := []*exec.Cmd{
			exec.Command("networksetup", "-setsocksfirewallproxy", "Wi-Fi", server, fmt.Sprintf("%d", port)),
			exec.Command("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "on"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil

	} else if runtime.GOOS == "windows" {
		proxyAddress := fmt.Sprintf("%s:%d", server, port)
		cmds := []*exec.Cmd{
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f"),
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/t", "REG_SZ", "/d", "socks="+proxyAddress, "/f"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil

	}

	return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

func ResetSocksProxy() error {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off")
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings",
			"/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	} else {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return cmd.Run()
}

func GetSocksProxyInfo() (*SocksProxyConfig, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("networksetup", "-getsocksfirewallproxy", "Wi-Fi")
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable")
	} else {
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// 执行命令并获取输出
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	output := out.String()

	// 解析 macOS/Linux 输出
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		scanner := bufio.NewScanner(strings.NewReader(output))
		ret := &SocksProxyConfig{}
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
				case "Server":
					ret.Server = value
					break
				case "Port":
					ret.Port, _ = strconv.Atoi(value)
					break
				}
			}
		}

		return ret, nil
	}

	// 解析 Windows 输出
	if runtime.GOOS == "windows" {
		enabled := strings.Contains(string(output), "0x1") // 0x1 代表开启代理

		cmd = exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err = cmd.Run(); err != nil {
			return nil, err
		}
		output := out.String()
		lines := strings.Split(string(output), "\n")
		proxyAddr := ""
		for _, line := range lines {
			if strings.Contains(line, "ProxyServer") {
				parts := strings.Fields(line)
				if len(parts) > 2 {
					proxyAddr = parts[len(parts)-1]
				}
			}
		}

		// 解析 SOCKS5 代理地址
		if strings.HasPrefix(proxyAddr, "socks=") {
			proxyAddr = strings.TrimPrefix(proxyAddr, "socks=")
		}

		// 解析 IP 和端口
		parts := strings.Split(strings.TrimSpace(proxyAddr), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid data:%s", proxyAddr)
		}

		ret := &SocksProxyConfig{}
		ret.Enable = enabled
		ret.Server = parts[0]
		if ret.Port, err = strconv.Atoi(parts[1]); err != nil {
			return nil, fmt.Errorf("invalid data:%s", proxyAddr)
		}

		return ret, nil
	}

	return nil, fmt.Errorf("failed to parse proxy settings")
}
