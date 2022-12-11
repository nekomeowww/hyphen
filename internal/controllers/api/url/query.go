package url

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/longkai/rfc7807"
	"github.com/nekomeowww/hyphen/pkg/handler"
)

type QueryURLParam struct {
	URL string `query:"url"`
}

type QueryURLResp struct {
	ShortURL string `json:"shortUrl"`
}

func (c *Controller) QueryURL(ctx echo.Context) (handler.Response, error) {
	var param QueryURLParam
	err := ctx.Bind(&param)
	if err != nil {
		return nil, rfc7807.New(rfc7807.InvalidArgument, fmt.Sprintf("invalid params: %v", err))
	}
	if param.URL == "" {
		return nil, rfc7807.New(rfc7807.InvalidArgument, "url must not be empty")
	}

	param.URL, err = c.URLs.NormalizeURL(param.URL)
	if err != nil {
		return nil, rfc7807.New(rfc7807.InvalidArgument, fmt.Sprintf("invalid url: %v", err))
	}

	result := c.URLs.FindOneShortURLByURL(param.URL)
	if result.IsError() {
		return nil, rfc7807.Wrap(rfc7807.Internal, "database error", result.Error())
	}
	if result.MustGet().ShortURL == "" {
		return nil, rfc7807.New(rfc7807.NotFound, "not found")
	}

	return &QueryURLResp{
		ShortURL: result.MustGet().ShortURL,
	}, nil
}

type QueryShortURLParam struct {
	Redirect bool   `query:"redirect"`
	URL      string `query:"url"`
}

type QueryShortURLResp struct {
	URL string `json:"url"`
}

func (c *Controller) QueryShortURL(ctx echo.Context) (handler.Response, error) {
	var param QueryShortURLParam
	err := ctx.Bind(&param)
	if err != nil {
		return nil, rfc7807.New(rfc7807.InvalidArgument, fmt.Sprintf("invalid params: %v", err))
	}

	param.URL, err = url.QueryUnescape(param.URL)
	if err != nil {
		return nil, rfc7807.New(rfc7807.InvalidArgument, fmt.Sprintf("invalid url: %v", err))
	}

	result := c.URLs.FindOneURLByShortURL(param.URL)
	if result.IsError() {
		return nil, rfc7807.Wrap(rfc7807.Internal, "database error", result.Error())
	}
	if result.MustGet().FullURL == "" {
		return nil, rfc7807.New(rfc7807.NotFound, "not found")
	}
	if param.Redirect {
		return nil, ctx.Redirect(http.StatusFound, result.MustGet().FullURL)
	}

	return &QueryShortURLResp{
		URL: result.MustGet().FullURL,
	}, nil
}
