# gin-zerolog ![maintenance-status](https://img.shields.io/badge/maintenance-deprecated-red.svg) [![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)
Use [`zerolog`](https://github.com/rs/zerolog)
within [`gin`](https://gin-gonic.com/docs/) applications.

_After I created this program I decided to make it part of a more general set of tools:_
[gin-utils](https://github.com/madkins23/gin-utils).

**This project is now deprecated.**

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

There are basic requirements when using `zerolog` within a `gin` application:

* provide a middleware function that writes records via `zerolog` and
* provide an IO writer object to replace the default `gin` logging stream, 
  trap the non-middleware log messages, and redirect them to `zerolog`.

[![Go Report Card](https://goreportcard.com/badge/github.com/madkins23/gin-zerolog)](https://goreportcard.com/report/github.com/madkins23/gin-zerolog)
![GitHub](https://img.shields.io/github/license/madkins23/gin-zerolog)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/madkins23/gin-zerolog)
[![Go Reference](https://pkg.go.dev/badge/github.com/madkins23/gin-zerolog.svg)](https://pkg.go.dev/github.com/madkins23/gin-zerolog)

# Usage

Import packages using:

    import (
        "github.com/gin-gonic/gin"
        "github.com/rs/zerolog"
        "github.com/rs/zerolog/log"

        "github.com/madkins23/gin-zerolog/ginzero"
    )

# Tools

There is a demo program located in `cmd/demo/demo.go`.

## Middleware

The basic logging for request traffic in `gin` is generally handled via middleware.
The existing default middleware sends request data to the default
logging streams with some formatting.

Add the `ginzero` logger using the following:

    router := gin.New() // not gin.Default()
    router.Use(ginzero.Logger())

Add routing configuration after these statements.

Use `gin.New()` instead of `gin.Default()`.
The latter adds its own logging middleware
which would conflict with the `ginzero` middleware.

## IO Writer

There is some `gin` logging of non-request issues that just goes to
the default logging streams.
This mostly happens at startup.
These streams can be replaced with any `IO.Writer` entity.

Trap and redirect these streams to `zerolog` using the following:

    gin.DefaultWriter = ginzero.NewWriter(zerolog.InfoLevel)
    gin.DefaultErrorWriter = ginzero.NewWriter(zerolog.ErrorLevel)
    router := gin.New() // or gin.Default() if not using ginzero.Logger()
