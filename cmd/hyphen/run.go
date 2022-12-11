package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nekomeowww/hyphen/internal/controllers"
	"github.com/nekomeowww/hyphen/internal/controllers/api"
	"github.com/nekomeowww/hyphen/internal/dao"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/internal/models"
	"github.com/nekomeowww/hyphen/internal/router"
	"github.com/nekomeowww/hyphen/pkg/meta"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start hyphen server, use -l and --data to specify listening address and data file path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`
_    _             _
| |  | |           | |
| |__| |_   _ _ __ | |__   ___ _ __
|  __  | | | | '_ \| '_ \ / _ \ '_ \
| |  | | |_| | |_) | | | |  __/ | | |
|_|  |_|\__, | .__/|_| |_|\___|_| |_|
         __/ | |
        |___/|_|

Hyphen URL shortener service (%s)`+"\n", meta.Version)

			app := fx.New(fx.Options(
				fx.Options(lib.NewModules()),
				fx.Provide(dao.NewBBolt(dbPath)),
				fx.Options(models.NewModules()),
				fx.Options(controllers.NewModules()),
				fx.Provide(router.NewRouter()),
				fx.Provide(api.NewAPI(listen)),
				fx.Invoke(api.Run()),
			))

			app.Run()
			stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), time.Second*15)
			defer stopCtxCancel()
			if err := app.Stop(stopCtx); err != nil {
				log.Fatal(err)
			}
		},
	}
)
