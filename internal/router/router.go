package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/internal/middlewares"
	"go.uber.org/fx"
)

type NewRouterParam struct {
	fx.In

	Logger *lib.Logger
}

type Router struct {
	Echo *echo.Echo
}

func NewRouter() func(NewRouterParam) *Router {
	return func(param NewRouterParam) *Router {
		e := echo.New()

		// logging
		e.Use(middlewares.LogRequest(param.Logger))
		// cores
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

		return &Router{
			Echo: e,
		}
	}
}
