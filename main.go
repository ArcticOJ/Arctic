//go:build ui || headless

package main

import (
	"context"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var Router = echo.New()

var OnDestroy []func()

var OnInit []func(context.Context)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	for i := range OnInit {
		OnInit[i](ctx)
	}
	go func() {
		<-ctx.Done()
		for i := range OnDestroy {
			OnDestroy[i]()
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
	logger.Panic(Router.Start(addr), "failed to listen on %s", addr)
}
