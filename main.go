//go:build ui || headless

package main

import (
	"blizzard/config"
	"blizzard/logger"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"os"
	"os/signal"
	"strings"
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
		Router.Use(middleware.BodyDump(func(c echo.Context, req, res []byte) {
			if strings.HasPrefix(c.Request().URL.Path, "/api") {
				logger.Logger.Debug().Str("url", c.Request().RequestURI).Bytes("req", req).Bytes("res", res).Msg("body")
			}
		}))
		Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Logger.Debug().
					Str("url", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Msg("req")
				return nil
			},
		}))
		rConf = middleware.RecoverConfig{
			DisableStackAll: true,
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				logger.Logger.Err(err).Str("url", c.Request().URL.RequestURI()).Send()
				fmt.Println(string(stack))
				return nil
			},
		}
	}
	Router.Use(middleware.RecoverWithConfig(rConf))
	Router.HideBanner = true
	Router.HidePort = true
	addr := net.JoinHostPort(config.Config.Host, fmt.Sprint(config.Config.Port))
	logger.Logger.Info().Msgf("starting server on %s", addr)
	logger.Logger.Fatal().Err(Router.Start(addr)).Send()
}
