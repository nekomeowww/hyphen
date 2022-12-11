package lib

import (
	"context"
	"log"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewLoggerParam struct {
	fx.In

	Lifecycle fx.Lifecycle
}

type Logger struct {
	*zap.Logger
}

func NewLogger() func(NewLoggerParam) (*Logger, error) {
	return func(param NewLoggerParam) (*Logger, error) {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				err = logger.Sync() // flushes buffer, if any
				if err != nil {
					log.Println(err)
				}
				return nil
			},
		})

		return &Logger{
			Logger: logger,
		}, nil
	}
}
