package mycrypto

import (
	"archive/zip"
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/myos"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// unzip 解压 ZIP 文件到指定目录
func Unzip(src string, dest string) error {
	// 打开 ZIP 文件
	mylog.Debug(myos.IsExistFile(src))

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("无法打开ZIP文件: %v", err)
	}
	defer func() {
		_ = zipReader.Close()
	}()

	// 遍历 ZIP 文件中的每个文件
	for _, file := range zipReader.File {
		filePath := filepath.Join(dest, file.Name)

		// 检查解压路径是否安全
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", filePath)
		}

		if file.FileInfo().IsDir() {
			// 创建目录
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("无法创建目录: %v", err)
			}
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("无法创建父目录: %v", err)
		}

		// 解压文件
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("无法创建文件: %v", err)
		}
		defer func() {
			_ = destFile.Close()
		}()

		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("无法打开文件: %v", err)
		}
		defer func() {
			_ = srcFile.Close()
		}()

		// 写入解压内容
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}

	return nil
}
