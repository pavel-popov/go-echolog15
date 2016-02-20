# Log15 middleware for Echo web-framework [![godoc reference](https://godoc.org/github.com/schmooser/go-echolog15?status.png)](https://godoc.org/github.com/schmooser/go-echolog15)


* [Echo web framework](https://labstack.com/echo/)
* [log15 logging framework]()

## Example

    package main

    import (
        "net/http"

        "github.com/labstack/echo"
        "github.com/schmooser/go-echolog15"
        log "gopkg.in/inconshreveable/log15.v2"
    )

    // Handler
    func hello(c *echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!\n")
    }

    func main() {
        // Echo instance
        e := echo.New()

        // Logger middleware
        e.Use(echolog15.Logger(log.New()))

        // Routes
        e.Get("/", hello)

        // Start server
        e.Run(":1323")
    }
