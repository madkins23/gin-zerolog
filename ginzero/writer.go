package ginzero

import (
	"fmt"
	"io"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Writer interface {
	io.Writer
}

func NewWriter(level zerolog.Level) Writer {
	return &writer{level: level}
}

type writer struct {
	level zerolog.Level
}

var (
	PTN_GIN_debug, _ = regexp.Compile("^\\s*\\[GIN-debug\\]\\s*")
	PTN_WARNING, _   = regexp.Compile("^\\s*\\[WARNING\\]\\s*")
)

func (w *writer) Write(p []byte) (n int, err error) {
	var sys string
	level := w.level

	// For the moment assume that a single Write() call is a single log record.
	msg := string(p)

	if match := PTN_GIN_debug.FindString(msg); match != "" {
		level = zerolog.DebugLevel
		msg = msg[len(match):]
		sys = "gin"
	}

	if match := PTN_WARNING.FindString(msg); match != "" {
		level = zerolog.WarnLevel
		msg = msg[len(match):]
	}

	var event *zerolog.Event
	switch level {
	case zerolog.DebugLevel:
		event = log.Debug()
	case zerolog.ErrorLevel:
		event = log.Error()
	case zerolog.InfoLevel:
		event = log.Info()
	case zerolog.WarnLevel:
		event = log.Warn()
	default:
		return 0, fmt.Errorf("unknown log level %s", w.level)
	}

	if sys != "" {
		event = event.Str("sys", sys)
	}
	event.Msg(msg)

	return len(p), nil
}
