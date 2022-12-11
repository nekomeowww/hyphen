package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/longkai/rfc7807"
)

type Response interface{}

type OkResponse struct {
	Data interface{} `json:"data"`
}

func NewHandler(handlerFunc func(echo.Context) (Response, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		r, err := handlerFunc(ctx)
		if err != nil {
			switch v := err.(type) {
			case *rfc7807.ProblemDetail:
				return ctx.JSON(v.Status, v)
			default:
				return ctx.JSON(http.StatusInternalServerError, rfc7807.Wrap(rfc7807.Internal, "unknown error occurred", v))
			}
		}

		return ctx.JSON(http.StatusOK, OkResponse{Data: r})
	}
}
