package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	logUtils "github.com/madkins23/go-utils/log"

	"github.com/madkins23/gin-zerolog/ginzero"
)

const appName = "template"

var port string

func main() {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	flags.StringVar(&port, "port", ":8080", "specify server port with leading colon")

	cof := logUtils.ConsoleOrFile{}
	cof.AddFlagsToSet(flags, "/tmp/console-or-file.log")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Printf("Error parsing command line flags: %s", err)
		}
		return
	}
	if err := cof.Setup(); err != nil {
		fmt.Printf("Log file creation error: %s", err)
		return
	}
	defer cof.CloseForDefer()

	// TODO: check port number

	gin.DefaultWriter = ginzero.NewWriter(zerolog.InfoLevel)
	gin.DefaultErrorWriter = ginzero.NewWriter(zerolog.ErrorLevel)
	router := gin.New() // not gin.Default()
	router.Use(ginzero.Logger())

	// Create context that listens for the interrupt signal from the OS.
	// NOTE: this assumes we're running on Linux.
	stopContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router.GET("/exit", func(c *gin.Context) {
		if proc, err := os.FindProcess(os.Getpid()); err != nil {
			log.Fatal().Err(err).Msg("Unable to find this process")
		} else if err = proc.Signal(syscall.SIGINT); err != nil {
			log.Fatal().Err(err).Msg("Unable send self interrupt signal")
		} else {
			log.Info().Msg("/exit invoked")
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "exiting",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Logger.Info().Msg(appName + " starting @ http://localhost" + port + "/ping")
	defer log.Logger.Info().Msg(appName + " finished")

	// Build http.Server object manually, don't use gin.Run().
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Start server in goroutine so shutdown code can run.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Running gin server")
		}
	}()

	// Listen for the interrupt signal.
	<-stopContext.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Msg("Shutting down gracefully, press Ctrl+C again to force exit.")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownContext); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
}
