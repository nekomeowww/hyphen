package url

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/longkai/rfc7807"
	"github.com/nekomeowww/hyphen/pkg/handler"
)

type NewParam struct {
	URL string `json:"url"`
}

type NewResp struct {
	URL      string `json:"url"`
	ShortURL string `json:"shortUrl"`
}

func (c *Controller) New(ctx echo.Context) (handler.Response, error) {
	var param NewParam
	err := ctx.Bind(&param)
	if err != nil {
		return nil, rfc7807.New(http.StatusBadRequest, fmt.Sprintf("invalid params: %v", err))
	}

	param.URL, err = c.URLs.NormalizeURL(param.URL)
	if err != nil {
		return nil, rfc7807.New(http.StatusBadRequest, fmt.Sprintf("invalid url: %v", err))
	}

	result := c.URLs.New(param.URL)
	if result.IsError() {
		return nil, rfc7807.Wrap(http.StatusInternalServerError, "datastore error", result.Error())
	}

	return &NewResp{
		URL:      param.URL,
		ShortURL: result.MustGet(),
	}, nil
}
