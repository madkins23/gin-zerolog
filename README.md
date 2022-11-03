# gin-zerolog
Use [`zerolog`](https://github.com/rs/zerolog)
within [`gin`](https://gin-gonic.com/docs/) applications.

There are basic requirements when using `zerolog` within a `gin` application:

* provide a middleware function that writes records via `zerolog` and
* provide an IO writer object to replace the default `gin` logging stream, 
  trap the non-middleware log messages, and redirect them to `zerolog`.

# Tools

## Middleware

The basic logging for request traffic in `gin` is generally handled via middleware.
The existing default middleware sends request data to the default
logging streams with some formatting.

Add the `ginzero` logger using the following:

    router := gin.Default() // or gin.New()
    router.Use(ginzero.Logger())

Add routing configuration after these statements.

## IO Writer

There is some `gin` logging of non-request issues that just goes to
the default logging streams.
This mostly happens at startup.
These streams can be replaced with any `IO.Writer` entity.

Trap and redirect these streams to `zerolog` using the following:

    gin.DefaultWriter = NewWriter(zerolog.InfoLevel)
    gin.DefaultErrorWriter = NewWriter(zerolog.ErrorLevel)
    router := gin.Default() // or gin.New()
