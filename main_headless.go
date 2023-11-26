//go:build headless

package main

import (
	"errors"
	"github.com/ArcticOJ/blizzard/v0"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/server"
	blizzardhttp "github.com/ArcticOJ/blizzard/v0/server/http"
	"github.com/ArcticOJ/blizzard/v0/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func init() {
	OnInit = append(OnInit, blizzard.Init)
	OnDestroy = append(OnDestroy, blizzard.Destroy)
	Router.IPExtractor = echo.ExtractIPFromXFFHeader()
	Router.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		code, message := http.StatusInternalServerError, "Internal Server Error"
		var er *echo.HTTPError
		if errors.As(err, &er) {
			code = er.Code
			message = er.Message.(string)
		}
		err = c.JSON(code, blizzardhttp.Error{Code: code, Message: message})
		if err != nil {
			Router.Logger.Error(err)
		}
	}
	Router.Validator = validator.New()
	if config.Config.Blizzard.EnableCORS {
		Router.Use(middleware.CORS())
	}
	server.Register(Router)
}
