package url

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/longkai/rfc7807"
	"github.com/nekomeowww/hyphen/pkg/handler"
)

type RevokeParam struct {
	ShortURL string `json:"shortUrl"`
}

type RevokeResp struct {
	Ok bool `json:"ok"`
}

func (c *Controller) Revoke(ctx echo.Context) (handler.Response, error) {
	var param RevokeParam
	err := ctx.Bind(&param)
	if err != nil {
		return nil, rfc7807.New(rfc7807.InvalidArgument, fmt.Sprintf("invalid params: %v", err))
	}

	result := c.URLs.RevokeOneShortURL(param.ShortURL)
	if result.IsError() {
		return nil, rfc7807.Wrap(rfc7807.Internal, "database error", result.Error())
	}

	return &RevokeResp{
		Ok: result.MustGet(),
	}, nil
}
