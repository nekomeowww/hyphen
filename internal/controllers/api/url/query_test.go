package url

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/nekomeowww/hyphen/pkg/handler"
	"github.com/nekomeowww/hyphen/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryURL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fullURL := fmt.Sprintf("https://example.com/你好世界/%s", utils.RandomHashString())
	normalizedURL, err := controller.URLs.NormalizeURL(fullURL)
	require.NoError(err)

	result := controller.URLs.New(normalizedURL)
	require.NoError(result.Error())
	assert.NotEmpty(result.MustGet())

	_, c := handler.NewTestContext(http.MethodGet, fmt.Sprintf("/url/full?url=%s", url.QueryEscape(fullURL)), nil)
	r, err := controller.QueryURL(c)
	require.NoError(err)
	require.NotNil(r)

	resp := r.(*QueryURLResp)
	assert.Equal(result.MustGet(), resp.ShortURL)
}

func TestQueryShortURL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fullURL := fmt.Sprintf("https://example.com/%s", utils.RandomHashString())
	normalizedURL, err := controller.URLs.NormalizeURL(fullURL)
	require.NoError(err)

	result := controller.URLs.New(normalizedURL)
	require.NoError(result.Error())
	assert.NotEmpty(result.MustGet())

	_, c := handler.NewTestContext(http.MethodGet, fmt.Sprintf("/url/short?url=%s", url.QueryEscape(result.MustGet())), nil)
	r, err := controller.QueryShortURL(c)
	require.NoError(err)
	require.NotNil(r)

	resp := r.(*QueryShortURLResp)
	assert.Equal(fullURL, resp.URL)
}
