package urls

import (
	"fmt"
	"os"
	"testing"

	"github.com/nekomeowww/hyphen/internal/dao"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type lifecycle struct{}

func (l *lifecycle) Append(hook fx.Hook) {}

var urls *URLModel

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

	urls = NewURLModel()(NewURLModelParam{
		Logger: logger,
		BBolt:  bbolt,
	})

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}

func TestURLs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var hash string
	fullURL := fmt.Sprintf("https://example.com/s/%s", utils.RandomHashString(64))

	t.Run("New", func(t *testing.T) {
		result := urls.New(fullURL)
		require.NoError(result.Error())

		hash = result.MustGet()
	})

	t.Run("Query", func(t *testing.T) {
		t.Run("URL", func(t *testing.T) {
			foundURLResult := urls.FindOneURLByShortURL(hash)
			require.NoError(foundURLResult.Error())

			foundURL := foundURLResult.MustGet()
			assert.Equal(fullURL, foundURL.FullURL)
		})
		t.Run("ShortURL", func(t *testing.T) {
			foundShortURLResult := urls.FindOneShortURLByURL(fullURL)
			require.NoError(foundShortURLResult.Error())

			foundShortURL := foundShortURLResult.MustGet()
			assert.Equal(hash, foundShortURL.ShortURL)
		})
	})
}

func TestRevokeShortURL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fullURL := fmt.Sprintf("https://example.com/s/%s", utils.RandomHashString(64))
	result := urls.New(fullURL)
	require.NoError(result.Error())

	hash := result.MustGet()
	revokeResult := urls.RevokeOneShortURL(hash)
	require.NoError(revokeResult.Error())
	assert.True(revokeResult.MustGet())

	foundURLResult := urls.FindOneURLByShortURL(hash)
	require.NoError(foundURLResult.Error())
	assert.Empty(foundURLResult.MustGet().FullURL)

	foundShortURLResult := urls.FindOneShortURLByURL(fullURL)
	require.NoError(foundShortURLResult.Error())
	assert.Empty(foundShortURLResult.MustGet().ShortURL)
}
