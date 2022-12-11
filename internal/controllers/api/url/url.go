package url

import (
	"github.com/labstack/echo/v4"
	"github.com/nekomeowww/hyphen/internal/dao"
	"github.com/nekomeowww/hyphen/internal/models/urls"
	"github.com/nekomeowww/hyphen/internal/router"
	"github.com/nekomeowww/hyphen/pkg/handler"
	"go.uber.org/fx"
)

type NewControllerParam struct {
	fx.In

	Router *router.Router
	BBolt  *dao.BBolt

	URLs *urls.URLModel
}

type Controller struct {
	Router *router.Router
	BBolt  *dao.BBolt

	URLs *urls.URLModel
}

func NewController() func(NewControllerParam) *Controller {
	return func(param NewControllerParam) *Controller {
		return &Controller{
			Router: param.Router,
			BBolt:  param.BBolt,
			URLs:   param.URLs,
		}
	}
}

func (c *Controller) Register(echo *echo.Group) {
	echo.GET("/url/full", handler.NewHandler(c.QueryURL))
	echo.GET("/url/short", handler.NewHandler(c.QueryShortURL))
	echo.POST("/url", handler.NewHandler(c.New))
	echo.DELETE("/url", handler.NewHandler(c.Revoke))
}
