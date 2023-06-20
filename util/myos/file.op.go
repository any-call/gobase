package myos

import (
	"os"
)

func IsExistPath(fPath string) bool {

	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsExistDir(fPath string) bool {
	info, err := os.Stat(fPath)
	if err != nil {
		return false
	}
	// 判断是否为目录
	if info.IsDir() {
		return true
	}

	return false
}

func IsExistFile(fPath string) bool {
	info, err := os.Stat(fPath)
	if err != nil {
		return false
	}
	// 判断是否为目录
	if info.IsDir() {
		return false
	}

	return true
}

func Remove(fpath string) error {
	// 删除文件
	err := os.Remove(fpath)
	if err != nil {
		// 如果删除失败，则尝试强制删除
		err = os.RemoveAll(fpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func Rename(oldPath, newPath string) error {
	// 如果新文件已经存在，则删除它
	_, err := os.Stat(newPath)
	if err == nil {
		err = os.Remove(newPath)
		if err != nil {
			return err
		}
	}
	// 重命名文件
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return err
}
