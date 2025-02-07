package myos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type (
	WebProxyConfig struct {
		Server string `json:"server"` // 服务器地址 (如 127.0.0.1)
		Port   int    `json:"port"`   // 端口号 (如 1080)
		Enable bool   `json:"enable"`
		//Auth     bool   // 是否需要身份验证 --目录命令设置不了
		//Username string // 用户名 (如果需要) --目录命令设置不了
		//Password string // 密码 (如果需要) --目录命令设置不了
	}
)

func (self WebProxyConfig) ToString() string {
	ret, _ := json.Marshal(self)
	return string(ret)
}

func SetWebProxy(server string, port int) error {
	// 设置 HTTP 代理
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// macOS 设置 HTTP 代理
		cmds := []*exec.Cmd{
			exec.Command("networksetup", "-setwebproxy", "Wi-Fi", server, fmt.Sprintf("%d", port)),
			exec.Command("networksetup", "-setwebproxystate", "Wi-Fi", "on"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	} else if runtime.GOOS == "windows" {
		// Windows 设置 HTTP 代理
		proxyServer := fmt.Sprintf("%s:%d", server, port)
		cmds := []*exec.Cmd{
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f"),
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/t", "REG_SZ", "/d", proxyServer, "/f"),
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

// 重置 HTTP 和 HTTPS 代理
func ResetWebProxy() error {
	// 重置 HTTP 代理
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmds := []*exec.Cmd{
			exec.Command("networksetup", "-setwebproxystate", "Wi-Fi", "off"),
		}

		for _, cmd := range cmds {
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		return nil
	} else if runtime.GOOS == "windows" {
		cmds := []*exec.Cmd{
			exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f"),
			exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/f"),
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

func GetWebProxyInfo() (*WebProxyConfig, error) {
	if runtime.GOOS == "windows" {
		// 获取代理开关
		cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		enabled := strings.Contains(string(output), "0x1")

		// 获取代理服务器地址
		cmd = exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer")
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}

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

		// 解析 IP 和端口
		if strings.HasPrefix(proxyAddr, "socks=") {
			proxyAddr = strings.TrimPrefix(proxyAddr, "socks=")
		}

		// 解析 IP 和端口
		parts := strings.Split(strings.TrimSpace(proxyAddr), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid data:%s", proxyAddr)
		}

		ret := WebProxyConfig{Enable: enabled}
		ret.Server = parts[0]
		if ret.Port, err = strconv.Atoi(parts[1]); err != nil {
			return nil, fmt.Errorf("invalid data:%s", proxyAddr)
		}

		return &ret, nil
	} else if runtime.GOOS == "darwin" {
		// macOS 使用 networksetup 获取 HTTP 代理
		cmd := exec.Command("networksetup", "-getwebproxy", "Wi-Fi")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		ret := &WebProxyConfig{}
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

	return nil, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}
