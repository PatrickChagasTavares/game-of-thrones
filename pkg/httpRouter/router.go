package httpRouter

import (
	"context"
	"net/http"
)

type (
	Router interface {
		Server(port string) error
		Get(path string, f HandlerFunc)
		Post(path string, f HandlerFunc)
		Put(path string, f HandlerFunc)
		Delete(paht string, f HandlerFunc)
		ParseHandler(h http.HandlerFunc) HandlerFunc
	}

	HandlerFunc func(ctx Context)

	Context interface {
		Context() context.Context
		JSON(statusCode int, data any)
		JSONError(err error)
		Decode(data any) error
		GetResponseWriter() http.ResponseWriter
		GetRequestReader() *http.Request
		GetQuery(param string) string
		GetParam(param string) string
		Validate(input any) error
	}
)
