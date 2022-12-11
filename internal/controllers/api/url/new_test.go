package url

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/nekomeowww/hyphen/pkg/handler"
	"github.com/nekomeowww/hyphen/pkg/types/dao/bbolt/keys"
	"github.com/nekomeowww/hyphen/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fullPath := fmt.Sprintf("https://example.com/%s", utils.RandomHashString())
	_, c := handler.NewTestContext(http.MethodPost, "/url", NewParam{
		URL: fullPath,
	})

	r, err := controller.New(c)
	require.NoError(err)

	resp := r.(*NewResp)
	assert.NotEmpty(resp.ShortURL)

	tx, err := controller.BBolt.Begin(false)
	require.NoError(err)
	defer func() {
		err = tx.Rollback()
		assert.NoError(err)
	}()

	bucket := tx.Bucket(keys.URLBucket.Format())
	require.NotNil(bucket)

	value := bucket.Get(keys.ShortenedURL1.Format(resp.ShortURL))
	require.NotNil(value)

	url, err := base64.StdEncoding.DecodeString(string(value))
	require.NoError(err)
	assert.Equal(resp.URL, string(url))
}
