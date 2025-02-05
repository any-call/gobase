package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func DockerExecCmd(containerName string) []string {
	return []string{"docker", "exec", containerName}
}

func BackupMySQL(dbUser, dbPassword, dbName, backupPath string, cmdModifier func() []string) (fullFile string, err error) {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s_%s.sql", backupPath, dbName, timestamp)

	// 基本命令
	cmdArgs := []string{
		"mysqldump",
		"-u", dbUser,
		"-p" + dbPassword,
		dbName,
	}

	// 如果 cmdModifier 存在，则插入其返回值
	if cmdModifier != nil {
		cmdArgs = append(cmdModifier(), cmdArgs...)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// 创建输出文件
	outFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = outFile.Close()
	}()

	// 设置输出到文件
	cmd.Stdout = outFile

	// 执行命令
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return fileName, nil
}
