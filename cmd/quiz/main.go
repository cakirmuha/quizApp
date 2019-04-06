package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"quizApp/service"
	"syscall"
	"time"

	"github.com/labstack/echo"
)

func main() {
	var (
		listenAddr = flag.String("listen", ":8181", "Listen addr")
		logLevel   = flag.String("log", "debug", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	ctx := context.Background()

	e, err := service.NewEcho(*logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sc, err := service.New(ctx, e.Logger, service.WithEcho(e), service.WithDB(true))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	cc := &ServerContext{Context: *sc}

	// Embed servercontext
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			copyCC := *cc
			copyCC.Context.Context = c

			return h(&copyCC)
		}
	})

	// Routes
	cc.SetupHandlers()

	go func() {
		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit
		cc.Log().Warnf("Received signal %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			panic(fmt.Sprintf("echo shutdown: %v", err))
		}
	}()

	// Start server
	if err := e.Start(*listenAddr); err != nil && err != http.ErrServerClosed {
		cc.Log().Error(err)
	} else {
		cc.Log().Infof("Shutting down... %v", err)
	}

}
