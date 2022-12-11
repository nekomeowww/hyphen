package url

import (
	"os"
	"testing"

	"github.com/nekomeowww/hyphen/internal/dao"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/internal/models/urls"
	"go.uber.org/fx"
)

type lifecycle struct{}

func (l *lifecycle) Append(hook fx.Hook) {}

var controller *Controller

func TestMain(m *testing.M) {
	bbolt, cleanupBBolt := dao.NewTestBBolt()
	defer func() {
		cleanupBBolt()
	}()

	logger, err := lib.NewLogger()(lib.NewLoggerParam{
		Lifecycle: &lifecycle{},
	})
	if err != nil {
		panic(err)
	}

	controller = &Controller{
		BBolt: bbolt,
		URLs: urls.NewURLModel()(urls.NewURLModelParam{
			Logger: logger,
			BBolt:  bbolt,
		}),
	}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}
