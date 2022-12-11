package models

import (
	"github.com/nekomeowww/hyphen/internal/models/urls"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(urls.NewURLModel()),
	)
}
