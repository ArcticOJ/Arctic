//go:build ui || headless

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var Router = echo.New()

var (
	// OnDestroy hooks are invoked after SIGINT or SIGTERM along with a 5-sec timeout as context, mainly for graceful shutdown of modules.
	OnDestroy []func(context.Context)
	// PostInit hooks are called after init() hooks.
	PostInit []func()
	// LateInit hooks are invoked even later than PostInit, used for HTTP server or less-crucial components' initialization
	LateInit []func(context.Context)
)

func StartServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	OnDestroy = append(OnDestroy, func(ctx context.Context) {
		logger.Panic(Router.Shutdown(ctx), "error shutting down http server")
	})
	defer cancel()
	for _, fn := range PostInit {
		fn()
	}
	for _, fn := range LateInit {
		fn(ctx)
	}
	go func() {
		<-ctx.Done()
		logger.Global.Info().Msg("server gracefully shutting down")
		_ctx, _cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer _cancel()
		for _, fn := range OnDestroy {
			fn(_ctx)
		}
	}()
	Router.Pre(middleware.RemoveTrailingSlash())
	rConf := middleware.RecoverConfig{
		DisablePrintStack: true,
		DisableStackAll:   true,
	}
	if config.Config.Debug {
		rConf = middleware.RecoverConfig{
			DisableStackAll: true,
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				logger.Global.Err(err).Str("url", c.Request().URL.RequestURI()).Msg("got panic")
				fmt.Println(string(stack))
				return nil
			},
		}
	}
	Router.Use(middleware.RecoverWithConfig(rConf))
	Router.HideBanner = true
	Router.HidePort = true
	addr := net.JoinHostPort(config.Config.Host, fmt.Sprint(config.Config.Port))
	logger.Global.Info().Msgf("server listening on %s", addr)
	if e := Router.Start(addr); !errors.Is(e, http.ErrServerClosed) {
		logger.Panic(e, "failed to listen on %s", addr)
	}
}
