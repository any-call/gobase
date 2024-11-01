package myos

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetSystemFontPath() ([]string, error) {
	// 获取当前操作系统
	osType := runtime.GOOS

	var fontDir string
	switch osType {
	case "windows":
		// Windows 系统常用字体目录
		fontDir = "C:\\Windows\\Fonts"
	case "linux":
		// Linux 系统常用字体目录
		fontDir = "/usr/share/fonts"
	case "darwin":
		// macOS 系统常用字体目录
		fontDir = "/Library/Fonts"
	default:
		return nil, fmt.Errorf("unsupported platform: %s", osType)
	}

	// 用于存储找到的所有字体文件路径
	var allFontFiles []string
	fontFiles, err := GetFontFiles(fontDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read font directory: %v", err)
	}
	allFontFiles = append(allFontFiles, fontFiles...)

	// 如果没有找到字体文件，返回一个提示信息
	if len(allFontFiles) == 0 {
		return nil, fmt.Errorf("no font files found in the specified directories")
	}

	return allFontFiles, nil
}

func GetFontFiles(dir string) ([]string, error) {
	var fontFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只保留 .ttf 和 .otf 文件
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".ttf" || filepath.Ext(info.Name()) == ".otf") {
			fontFiles = append(fontFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fontFiles, nil
}
