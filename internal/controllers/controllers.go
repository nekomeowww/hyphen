package controllers

import (
	"github.com/nekomeowww/hyphen/internal/controllers/api/url"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(url.NewController()),
	)
}
