package mylog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	FmtEmptySeparate = ""
)

type Formatter interface {
	Format(entry *Entry) error
}

type TextFormatter struct {
	IgnoreBasicFields bool
}

func (f *TextFormatter) Format(e *Entry) error {
	var content string
	if !f.IgnoreBasicFields {
		content += fmt.Sprintf("%s %s", e.Time.Format(time.RFC3339), levelNameMapping[e.Level]) // allocs
		if e.File != "" {
			short := e.File
			for i := len(e.File) - 1; i > 0; i-- {
				if e.File[i] == '/' {
					short = e.File[i+1:]
					break
				}
			}
			content += fmt.Sprintf(" %s:%d", short, e.Line)
		}
		e.Buffer.WriteString(" ")
	}

	switch e.Format {
	case FmtEmptySeparate:
		content += fmt.Sprint(e.Args...)
	default:
		content += fmt.Sprintf(e.Format, e.Args...)
	}

	if style, ok := levelStyleMapping[e.Level]; ok {
		content = style(content)
	}

	e.Buffer.WriteString(content)
	e.Buffer.WriteString("\n")
	return nil
}

type JsonFormatter struct {
	IgnoreBasicFields bool
}

func (f *JsonFormatter) Format(e *Entry) error {
	if !f.IgnoreBasicFields {
		e.Map["level"] = levelNameMapping[e.Level]
		e.Map["time"] = e.Time.Format(time.RFC3339)
		if e.File != "" {
			e.Map["file"] = e.File + ":" + strconv.Itoa(e.Line)
			e.Map["func"] = e.Func
		}

		switch e.Format {
		case FmtEmptySeparate:
			e.Map["message"] = fmt.Sprint(e.Args...)
			break
		default:
			e.Map["message"] = fmt.Sprintf(e.Format, e.Args...)
			break
		}

		b, err := json.Marshal(e.Map)
		if err != nil {
			return err
		}

		var content string
		if style, ok := levelStyleMapping[e.Level]; ok {
			content = style(string(b))
		} else {
			content = string(b)
		}

		e.Buffer.WriteString(content)
		return nil
	}

	var content string
	switch e.Format {
	case FmtEmptySeparate:
		for _, arg := range e.Args {
			b, err := json.Marshal(arg)
			if err != nil {
				return err
			}
			content += string(b)
		}
		break

	default:
		content = fmt.Sprintf(e.Format, e.Args...)
		break
	}

	if style, ok := levelStyleMapping[e.Level]; ok {
		content = style(content)
	}

	e.Buffer.WriteString(content)
	return nil
}
