package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nekomeowww/hyphen/internal/lib"
	"go.uber.org/zap"
)

func LogRequest(logger *lib.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("uri", values.URI),
				zap.String("host", values.Host),
				zap.String("remote_ip", values.RemoteIP),
				zap.Int64("response_time", int64(values.Latency.Milliseconds())),
				zap.Int64("status", int64(values.Status)),
			)

			return nil
		},
	})
}
