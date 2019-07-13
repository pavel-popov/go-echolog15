package echolog15

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/labstack/echo"
)

// LogProvider provides function required for logging.
type LogProvider interface {
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

// Logger is a logger middleware for log15 package.
func Logger(l LogProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			l.Debug("Echo request", "request", func() string {
				dump, _ := httputil.DumpRequest(req, true)
				return string(dump)
			}())

			remoteAddr := req.RemoteAddr
			if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			l.Info("Echo response",
				"remoteAddr", remoteAddr,
				"method", req.Method,
				"path", path,
				"status", res.Status,
				"time", stop.Sub(start),
				"size", res.Size)
			return nil
		}
	}
}

// HTTPErrorHandler is an error handler with log15 support.
func HTTPErrorHandler(l LogProvider) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message.(string)
		}
		if !c.Response().Committed {
			c.String(code, msg)
		}
		request := func() string {
			if c.Request().Method == "GET" {
				dump, _ := httputil.DumpRequest(c.Request(), true)
				return string(dump)
			}
			return "Request body dumped only for GET requests"
		}()
		if err.Error() != msg {
			l.Error("Echo error", "err", err, "code", code, "msg", msg, "request", request)
		} else {
			l.Error("Echo error", "code", code, "msg", msg, "request", request)
		}
	}
}
