/*
demo of ginzero package to support using [zerolog] with [gin].

Once running the application responds to http://:55555/ping with:
  {"message":"pong"}

Usage:
  demo

[gin]: https://gin-gonic.com/docs/
[zerolog]: https://github.com/rs/zerolog
*/
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/madkins23/gin-zerolog/ginzero"
)

func main() {
	gin.DefaultWriter = ginzero.NewWriter(zerolog.InfoLevel)
	gin.DefaultErrorWriter = ginzero.NewWriter(zerolog.ErrorLevel)
	router := gin.New()
	router.Use(ginzero.Logger())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	if err := router.Run(":55555"); err != nil {
		log.Error().Err(err).Msg("Server failure")
	}
}
