package myos

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func IsExistUrlFile(url string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Check the HTTP status code
	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

func Filename(fpath string) string {
	filename := filepath.Base(fpath)
	return filename
}

func Dir(fpath string) string {
	dir := filepath.Dir(fpath)
	return dir
}

func FindFilesWithExtRecursive(dir string, extension string) ([]string, error) {
	var files []string

	// 遍历目录及子目录
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为文件，并且是否匹配其中一个指定的扩展名
		if !d.IsDir() && strings.HasSuffix(d.Name(), extension) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func FindFilesWithExt(dir string, extension string) ([]string, error) {
	var filesWithExtension []string

	// 遍历文件夹
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是文件且扩展名匹配
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			filesWithExtension = append(filesWithExtension, path)
		}
		return nil
	})

	return filesWithExtension, err
}
