package service

import (
	"errors"
	"fmt"
	"net/http"
	"quizApp/model"
	"strings"
	"time"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

const (
	allowedHostHeader = "fasttract.io" // TODO: block other hosts
	routeNameCtxVar   = "routeName"
)

func NewEcho(logLevel string) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true

	switch strings.ToLower(logLevel) {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "warn":
		e.Logger.SetLevel(log.WARN)
	case "error":
		e.Logger.SetLevel(log.ERROR)
	default:
		return nil, fmt.Errorf("Invalid log level specified")
	}

	// Recover from panics
	e.Use(middleware.Recover())

	// Configure logger
	{
		lcfg := DefaultLoggerConfig
		lcfg.Skipper = func(c echo.Context) bool {
			return c.Request().RequestURI == "/" || c.Request().Method == http.MethodOptions
		}
		e.Use(LoggerWithConfig(lcfg))
	}

	// Middleware to ignore requests from different hosts
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !strings.HasSuffix(c.Request().Host, allowedHostHeader) && !strings.HasPrefix(c.Request().Host, "localhost") {
				return model.NewApiError(model.ApiErrorNotFound, "", nil)
			}
			return h(c)
		}
	})

	return e, nil
}

func (cc *Context) SetupHandlers(prefix string) (*echo.Group, func()) {
	e := cc.Echo()

	e.GET("/", func(c echo.Context) error {
		return nil
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		return nil
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	routes := make(map[string]string)

	api := e.Group(prefix, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			routeName := routes[c.Request().Method+":"+c.Path()]
			c.Set(routeNameCtxVar, routeName)
			return next(c)
		}
	}, ApiErrorMiddleware)

	return api, func() {
		for _, r := range e.Routes() {
			routes[r.Method+":"+r.Path] = r.Name
		}
	}
}

// ApiErrorMiddleware logs the internal error and returns a safe error code/message to the client, using ApiResponse
func ApiErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		// Convert our non-ApiResponse errors to internal server error...
		if _, ok := err.(*model.ApiResponseError); !ok {
			err = model.NewApiError(model.ApiErrorInternal, "", err)
		}

		code := 0
		msg := ""
		e := err // "internal" error to log...

		if he, ok := e.(*model.ApiResponseError); ok {
			code = he.Code
			msg = fmt.Sprintf("%v", he.Message)
			e = he.Internal()
		}

		if he, ok := e.(*echo.HTTPError); ok {
			if he.Code == http.StatusNotFound {
				return e
			}

			if code == 0 {
				code = he.Code
			}

			emsg := fmt.Sprintf("%v", he.Message)
			if msg == "" {
				msg = emsg
			}
			if he.Internal != nil {
				e = he.Internal
			} else {
				e = errors.New(emsg)
			}
		}

		if code == 0 {
			code = http.StatusInternalServerError
		}
		if msg == "" {
			msg = err.Error()
		}

		m := map[string]interface{}{
			"id":    c.Response().Header().Get(echo.HeaderXRequestID),
			"error": msg,
		}
		if e != nil && e != err {
			m["internal"] = e.Error()
		}

		if code != http.StatusNotFound {
			c.Logger().Errorj(m)
		}

		return ApiResponse(c, code, model.NewApiError(code, msg, nil))
	}
}

// NoCacheMiddleware removes cache headers.
// "Inspired" from https://github.com/zenazn/goji/blob/master/web/middleware/nocache.go
func NoCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// Taken from https://github.com/mytrile/nocache
	var noCacheHeaders = map[string]string{
		"Expires":         time.Unix(0, 0).Format(time.RFC1123),
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	return func(c echo.Context) error {
		// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
		//      Expires: Thu, 01 Jan 1970 00:00:00 UTC
		//      Cache-Control: no-cache, private, max-age=0
		//      X-Accel-Expires: 0
		//      Pragma: no-cache (for HTTP/1.0 proxies/clients)

		r := c.Response()

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			r.Header().Set(k, v)
		}

		return next(c)
	}
}

func EnsureContentTypeFunc(contentType string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ct := strings.Split(c.Request().Header.Get(echo.HeaderContentType), ";")[0]
			if ct != contentType {
				msg := fmt.Sprintf("Invalid content-type %q, expecting %q", ct, contentType)
				return model.NewApiError(model.ApiErrorBadRequest, msg, nil)
			}
			return next(c)
		}
	}
}
