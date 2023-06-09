package mylog

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var levelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

var levelStyleMapping = map[Level]func(string) string{
	DebugLevel: debugStyle,
	InfoLevel:  infoStyle,
	WarnLevel:  warnStyle,
	ErrorLevel: errorStyle,
	PanicLevel: panicStyle,
	FatalLevel: fatalStyle,
}

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

type Level uint8

func (l *Level) unmarshalText(text []byte) bool {
	if text == nil {
		return false
	}

	switch strings.ToLower(string(text)) {
	case "debug":
		*l = DebugLevel
		break

	case "info":
		*l = InfoLevel
		break

	case "warn":
		*l = WarnLevel
		break

	case "error":
		*l = ErrorLevel
		break

	case "panic":
		*l = PanicLevel
		break

	case "fatal":
		*l = FatalLevel
		break

	default:
		return false
	}

	return true
}

func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) {
		return fmt.Errorf("unrecognized level: %q", text)
	}

	return nil
}
