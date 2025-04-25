package mylog

import (
	"io"
	"os"
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
	opt       *options
	mu        sync.Mutex
	entryPool *sync.Pool
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

func (l *logger) Debug(args ...interface{}) {
	l.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.entry().write(PanicLevel, FmtEmptySeparate, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.entry().write(FatalLevel, FmtEmptySeparate, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry().write(DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.entry().write(InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry().write(WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry().write(ErrorLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.entry().write(PanicLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
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

// std logger
func Debug(args ...interface{}) {
	std.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func Info(args ...interface{}) {
	std.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func Warn(args ...interface{}) {
	std.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func Error(args ...interface{}) {
	std.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func Panic(args ...interface{}) {
	std.entry().write(PanicLevel, FmtEmptySeparate, args...)
}

func Fatal(args ...interface{}) {
	std.entry().write(FatalLevel, FmtEmptySeparate, args...)
}

func Debugf(format string, args ...interface{}) {
	std.entry().write(DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	std.entry().write(InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.entry().write(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.entry().write(ErrorLevel, format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.entry().write(PanicLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	std.entry().write(FatalLevel, format, args...)
}
