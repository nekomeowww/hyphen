package router

import (
	"github.com/labstack/echo/v4"
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
		e.Use(middlewares.LogRequest(param.Logger))
		return &Router{
			Echo: e,
		}
	}
}
