//go:build headless

package main

import (
	"blizzard"
	"blizzard/config"
	"blizzard/models"
	blizzardhttp "blizzard/server/http"
	"blizzard/server/http/middlewares"
	"blizzard/validator"
	"errors"
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
		err = c.JSON(code, models.Error{Code: code, Message: message})
		if err != nil {
			Router.Logger.Error(err)
		}
	}
	Router.Validator = validator.New()
	if config.Config.EnableCORS {
		Router.Use(middleware.CORS())
	}
	if config.Config.RateLimit > 0 {
		Router.Use(middlewares.RateLimit())
	}
	blizzardhttp.Register(Router)
}
