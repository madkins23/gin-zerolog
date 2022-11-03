package ginzero

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Writer interface for replacing gin standard output and/or error streams.
type Writer interface {
	io.Writer
}

// NewWriter returns a Writer object with the specified zerolog.Level.
// There are two gin output streams: gin.DefaultWriter and gin.DefaultErrorWriter.
// These streams are used by gin internal code outside the request middleware loop.
// Create a separate Writer object with a different zerolog.Level for each stream
// or create a single object for both streams (untested but should work).
func NewWriter(level zerolog.Level) Writer {
	return &writer{level: level}
}

// Make sure the writer struct implements ginzero.Writer.
var _ = Writer(&writer{})

// writer object returned by NewWriter function.
type writer struct {
	// Default zerolog level for this object.
	// Can be overridden by error levels (specified in logLevels variable)
	// in square brackets at the beginning of a log record line.
	level zerolog.Level
}

var (
	logLevels = map[string]zerolog.Level{
		"DEBUG":   zerolog.DebugLevel,
		"ERROR":   zerolog.ErrorLevel,
		"INFO":    zerolog.InfoLevel,
		"WARNING": zerolog.WarnLevel,
	}
	ptn_GIN, _       = regexp.Compile("^\\s*\\[GIN\\]\\s*")
	ptn_GIN_debug, _ = regexp.Compile("^\\s*\\[GIN-debug\\]\\s*")
	ptn_log_level, _ = regexp.Compile("^\\s*\\[(DEBUG|ERROR|INFO|WARNING)\\]\\s*")
)

// Write a block of data to the (supposedly) stream object.
// For the moment we're assuming that there is a single Write() call for each log record.
// TODO: Fix code to handle multiple Write() calls per log record.
func (w *writer) Write(p []byte) (n int, err error) {
	level := w.level
	msg := strings.TrimRight(string(p), "\n")
	var sys string

	for x := 0; x < 10; x++ { // Don't use infinite for loop for safety
		// Pull off prefix sequences that represent log information.
		if match := ptn_GIN.FindString(msg); match != "" {
			msg = msg[len(match):]
			sys = "gin"
		} else if match := ptn_GIN_debug.FindString(msg); match != "" {
			level = zerolog.DebugLevel
			msg = msg[len(match):]
			sys = "gin"
		} else if matches := ptn_log_level.FindStringSubmatch(msg); len(matches) > 1 {
			var ok bool
			if level, ok = logLevels[matches[1]]; !ok {
				return 0, fmt.Errorf("no level %s", matches[1])
			}
			msg = msg[len(matches[0]):]
		} else {
			break
		}
	}

	// Create the initial zerolog.Event object with the specified level.
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
