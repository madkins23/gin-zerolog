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

## IO Writer

There is a log of `gin` logging of non-request issues that just goes to
the default logging streams.
These streams can be replaced with any `IO.Writer` entity.
