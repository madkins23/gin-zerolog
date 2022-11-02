package ginzero

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns a Gin middleware function that generates a zerolog record for the current request.
// The record will be generated in the format for which zerolog has been configured.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		var event *zerolog.Event
		code := c.Writer.Status()
		if code >= 400 && code < 500 {
			event = log.Warn()
		} else if code >= 500 {
			event = log.Error()
		} else {
			event = log.Debug()
		}
		event = event.Int("code", code)

		event = event.Time("time", start)
		event = event.Dur("dur", duration)
		event = event.Str("ip", c.ClientIP())

		event = event.Str("meth", c.Request.Method)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		event = event.Str("path", path)

		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}
		event.Msg(msg)
	}
}
