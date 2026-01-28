package mylog

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var std = New()

func New(opts ...Option) *logger {
	logger := &logger{opt: initOptions(opts...)}
	logger.entryPool = &sync.Pool{New: func() interface{} { return entry(logger) }}
	return logger
}

type logger struct {
	opt        *options
	mu         sync.Mutex
	entryPool  *sync.Pool
	disableLog bool //写日志开关
}

func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, opt := range opts {
		opt(l.opt)
	}
}

func (l *logger) entry() *Entry {
	return l.entryPool.Get().(*Entry)
}

func (l *logger) Write(data []byte) (int, error) {
	l.entry().write(l.opt.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
	return 0, nil
}

func (l *logger) Writer() io.Writer {
	return l
}

func (l *logger) SetDisableLog(b bool) {
	l.disableLog = b
}

func (l *logger) Debug(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func (l *logger) Info(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func (l *logger) Warn(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func (l *logger) Error(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func (l *logger) Panic(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(PanicLevel, FmtEmptySeparate, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(FatalLevel, FmtEmptySeparate, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.disableLog {
		return
	}

	l.entry().write(InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(ErrorLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(PanicLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.disableLog {
		return
	}
	l.entry().write(FatalLevel, format, args...)
}

func (l *logger) InitLevel(level Level, en bool) *logger {
	l.SetOptions(WithLevel(level, en))
	return l
}

func StdLogger() *logger {
	return std
}

func SetOptions(opts ...Option) {
	std.SetOptions(opts...)
}

func GenFileByFilePath(fullPathFile string) (*os.File, error) {
	_, err := os.Stat(fullPathFile)
	if err != nil {
		if err := os.MkdirAll(fullPathFile, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(fullPathFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func InitDebugLevel(en bool) *logger {
	return std.InitLevel(DebugLevel, en)
}

func InitInfoLevel(en bool) *logger {
	return std.InitLevel(InfoLevel, en)
}

func InitWarnLevel(en bool) *logger {
	return std.InitLevel(WarnLevel, en)
}

func InitErrorLevel(en bool) *logger {
	return std.InitLevel(ErrorLevel, en)
}

func InitPanicLevel(en bool) *logger {
	return std.InitLevel(PanicLevel, en)
}

func InitFatalLevel(en bool) *logger {
	return std.InitLevel(FatalLevel, en)
}

func DisableLog(b bool) {
	std.SetDisableLog(b)
}

// std logger
func DebugKV(event string, kv ...any) {
	std.Debug(buildLogfmt(event, kv...))
}

func InfoKV(event string, kv ...any) {
	std.Info(buildLogfmt(event, kv...))
}

func WarnKV(event string, kv ...any) {
	std.Warn(buildLogfmt(event, kv...))
}

func ErrorKV(event string, kv ...any) {
	std.Error(buildLogfmt(event, kv...))
}

func PanicKV(event string, kv ...any) {
	std.Panic(buildLogfmt(event, kv...))
}

func FatalKV(event string, kv ...any) {
	std.Fatal(buildLogfmt(event, kv...))
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func Panic(args ...interface{}) {
	std.Panic(args...)
}

func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// kv 日志格式
func buildLogfmt(event string, kv ...any) string {
	var b strings.Builder
	b.WriteString("event=")
	b.WriteString(event)

	for i := 0; i < len(kv); i += 2 {
		if i+1 >= len(kv) {
			break
		}

		key := fmt.Sprint(kv[i])
		val := fmt.Sprint(kv[i+1])

		// 防止 value 中有空格
		if strings.ContainsAny(val, " \t") {
			val = strconv.Quote(val)
		}

		b.WriteString(" ")
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(val)
	}

	return b.String()
}
