package myfile

import (
	"os"
	"path/filepath"
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDir(path string) bool {
	is, _ := IsDirE(path)
	return is
}

func IsFile(path string) bool {
	is, _ := IsFileE(path)
	return is
}

func IsDirE(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true, nil
	}
	return false, err
}

func IsFileE(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.Mode().IsRegular() {
		return true, nil
	}
	return false, err
}

func IsSymlink(path string) bool {
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}

func IsSymlinkE(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		return true, nil
	}
	return false, err
}

// IsShortcut checks a file whether is a shortcut on Windows.
func IsShortcut(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".lnk" {
		return true
	}
	return false
}
