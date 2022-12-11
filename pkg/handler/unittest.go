package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/labstack/echo/v4"
)

func NewTestContext(method, endpoint string, body interface{}) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()

	var req *http.Request
	switch t := body.(type) {
	case nil:
		req, _ = http.NewRequest(method, endpoint, nil)
	case io.Reader:
		req, _ = http.NewRequest(method, endpoint, t)
		req.Header.Set("Content-Type", "application/json")
	case string:
		query, err := url.ParseQuery(t)
		if err != nil {
			panic(err)
		}

		req, _ = http.NewRequest(method, endpoint, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.PostForm = query
	case url.Values:
		req, _ = http.NewRequest(method, endpoint, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.PostForm = t
	default:
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		req, _ = http.NewRequest(method, endpoint, bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return rec, ctx
}
