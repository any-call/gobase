package mylog

import (
	"io"
	"os"
)

// log options
type options struct {
	output        io.Writer
	level         map[Level]bool
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
	disableColors bool
}

type Option func(*options)

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = os.Stderr
	}

	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	if o.level == nil {
		o.level = make(map[Level]bool, 5)
		o.level[DebugLevel] = true
		o.level[InfoLevel] = true
		o.level[WarnLevel] = true
		o.level[ErrorLevel] = true
		o.level[PanicLevel] = true
		o.level[FatalLevel] = true
	}

	return
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithDisableColor(en bool) Option {
	return func(o *options) {
		o.disableColors = en
	}
}

func WithLevel(level Level, en bool) Option {
	return func(o *options) {
		o.level[level] = en
	}
}

func WithStdLevel(level Level) Option {
	return func(o *options) {
		o.stdLevel = level
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller
	}
}
