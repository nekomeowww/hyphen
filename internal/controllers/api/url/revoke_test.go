package url

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nekomeowww/hyphen/pkg/handler"
	"github.com/nekomeowww/hyphen/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevoke(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fullURL := fmt.Sprintf("https://example.com/%s", utils.RandomHashString())
	result := controller.URLs.New(fullURL)
	require.NoError(result.Error())
	assert.NotEmpty(result.MustGet())

	_, c := handler.NewTestContext(http.MethodGet, "/url", RevokeParam{
		ShortURL: result.MustGet(),
	})

	r, err := controller.Revoke(c)
	require.NoError(err)
	require.NotNil(r)

	findResult := controller.URLs.FindOneURLByShortURL(result.MustGet())
	require.NoError(findResult.Error())
	assert.Empty(findResult.MustGet().FullURL)
}
