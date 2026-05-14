package mylog

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const fileTimeFormat = "2006-01-02 15:04:05.000"

// Field 是一条文件日志中的扩展字段。
type Field struct {
	Key   string
	Value any
}

// FileLogger 提供按文件写入与读取日志的能力。
type FileLogger struct {
	mu  sync.Mutex
	dir string
}

var stdFileLogger = &FileLogger{dir: "logs"}

// F 创建一个日志扩展字段。
func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// NewFileLogger 创建一个使用指定目录的文件日志器。
func NewFileLogger(dir string) *FileLogger {
	if dir == "" {
		dir = "logs"
	}
	return &FileLogger{dir: dir}
}

// SetDir 设置文件日志目录。
func (l *FileLogger) SetDir(dir string) {
	if dir == "" {
		dir = "logs"
	}

	l.mu.Lock()
	l.dir = dir
	l.mu.Unlock()
}

// Write 向指定日志文件写入一条 JSON line 日志。
func (l *FileLogger) Write(filename, level, content string, fields ...Field) {
	if filename == "" {
		return
	}

	entry := map[string]any{
		"timestamp": time.Now().Format(fileTimeFormat),
		"level":     level,
		"content":   content,
		"caller":    shortCaller(3),
	}
	for _, field := range fields {
		if field.Key != "" {
			entry[field.Key] = field.Value
		}
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err = os.MkdirAll(l.dir, 0755); err != nil {
		return
	}

	file, err := os.OpenFile(filepath.Join(l.dir, filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = file.Write(append(data, '\n'))
}

// ReadLastNLines 读取当前日志目录下指定文件的最后 n 行。
func (l *FileLogger) ReadLastNLines(filename string, n int) ([]string, error) {
	l.mu.Lock()
	path := filepath.Join(l.dir, filename)
	l.mu.Unlock()

	return ReadLastNLines(path, n)
}

// KeepLastNLines 裁剪当前日志目录下指定文件，只保留最后 n 行。
func (l *FileLogger) KeepLastNLines(filename string, n int) error {
	l.mu.Lock()
	path := filepath.Join(l.dir, filename)
	l.mu.Unlock()

	return KeepLastNLines(path, n)
}

// SetLogDir 设置默认文件日志目录。
func SetLogDir(dir string) {
	stdFileLogger.SetDir(dir)
}

// WriteLogFile 使用默认文件日志器写入一条 JSON line 日志。
func WriteLogFile(filename string, level string, content string, fields ...Field) {
	stdFileLogger.Write(filename, level, content, fields...)
}

// ReadLastNLogLines 读取默认日志目录下指定文件的最后 n 行。
func ReadLastNLogLines(filename string, n int) ([]string, error) {
	return stdFileLogger.ReadLastNLines(filename, n)
}

// KeepLastNLogLines 裁剪默认日志目录下指定文件，只保留最后 n 行。
func KeepLastNLogLines(filename string, n int) error {
	return stdFileLogger.KeepLastNLines(filename, n)
}

// ReadLastNLines 读取指定文件路径的最后 n 行。
func ReadLastNLines(path string, n int) ([]string, error) {
	if n <= 0 {
		return []string{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		return []string{}, nil
	}

	bufferSize := int64(4096)
	if bufferSize > info.Size() {
		bufferSize = info.Size()
	}

	buffer := make([]byte, bufferSize)
	position := info.Size()
	lineCount := 0

	for lineCount < n && position > 0 {
		readSize := bufferSize
		if position < bufferSize {
			readSize = position
		}

		position -= readSize
		if _, err = file.Seek(position, io.SeekStart); err != nil {
			return nil, err
		}
		if _, err = file.Read(buffer[:readSize]); err != nil {
			return nil, err
		}

		for i := readSize - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				lineCount++
				if lineCount > n {
					position += int64(i) + 1
					break
				}
			}
		}
	}

	if _, err = file.Seek(position, io.SeekStart); err != nil {
		return nil, err
	}

	lines := make([]string, 0, n)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

// KeepLastNLines 裁剪指定文件，只保留最后 n 行。
func KeepLastNLines(path string, n int) error {
	if n < 0 {
		n = 0
	}

	lines, err := ReadLastNLines(path, n)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err = file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func shortCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	file = filepath.ToSlash(file)
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}

	return file + ":" + strconv.Itoa(line)
}
