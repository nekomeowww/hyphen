package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/nekomeowww/hyphen/internal/controllers/api/url"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/internal/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewAPIParam struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *lib.Logger
	Router    *router.Router
	URL       *url.Controller
}

type API struct {
	*http.Server
}

func NewAPI(addr string) func(NewAPIParam) *API {
	return func(param NewAPIParam) *API {
		apiGroup := param.Router.Echo.Group("/api")
		v1Group := apiGroup.Group("/v1")

		param.URL.Register(v1Group)

		for _, v := range param.Router.Echo.Routes() {
			param.Logger.Debug("registered route", zap.String("method", v.Method), zap.String("path", v.Path))
		}

		server := &http.Server{
			Addr:    addr,
			Handler: param.Router.Echo,
		}

		param.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := server.Shutdown(closeCtx); err != nil && err != http.ErrServerClosed {
					param.Logger.Error("shutdown server failed", zap.Error(err))
					return err
				}

				return nil
			},
		})

		api := &API{Server: server}
		return api
	}
}

func Run() func(logger *lib.Logger, api *API) error {
	return func(logger *lib.Logger, api *API) error {
		logger.Info("starting api server...")
		listener, err := net.Listen("tcp", api.Addr)
		if err != nil {
			return fmt.Errorf("failed to listen %s: %v", api.Addr, err)
		}

		go func() {
			err = api.Serve(listener)
			if err != nil && err != http.ErrServerClosed {
				logger.Fatal(err.Error())
			}
		}()

		logger.Info("api server listening...", zap.String("addr", api.Addr))
		return nil
	}
}
