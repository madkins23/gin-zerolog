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
	PTN_LEVEL, _     = regexp.Compile("^\\s*\\[(DEBUG|ERROR|INFO|WARNING)\\]\\s*")
)

var (
	levels = map[string]zerolog.Level{
		"DEBUG":   zerolog.DebugLevel,
		"ERROR":   zerolog.ErrorLevel,
		"INFO":    zerolog.InfoLevel,
		"WARNING": zerolog.WarnLevel,
	}
)

func (w *writer) Write(p []byte) (n int, err error) {
	var sys string
	level := w.level

	// For the moment assume that a single Write() call is a single log record.
	msg := string(p)

	for x := 0; x < 10; x++ { // Don't use infinite for loop for safety
		// Pull off prefix sequences that represent log information.
		if match := PTN_GIN_debug.FindString(msg); match != "" {
			level = zerolog.DebugLevel
			msg = msg[len(match):]
			sys = "gin"
		} else if matches := PTN_LEVEL.FindStringSubmatch(msg); len(matches) > 1 {
			var ok bool
			if level, ok = levels[matches[1]]; !ok {
				return 0, fmt.Errorf("no level %s", matches[1])
			}
			msg = msg[len(matches[0]):]
		} else {
			break
		}
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
