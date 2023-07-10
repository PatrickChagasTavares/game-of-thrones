package httpRouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	ginRouter struct {
		router *gin.Engine
	}
)

const (
	ginTracerKey = "gin-tracer"
)

func NewGinRouter() Router {
	router := gin.Default()

	router.Use(
		// Set the content type default = application/json
		setContentType("application/json"),
		// Set middleware to tracer
		tracing(),
	)

	return &ginRouter{
		router: router,
	}
}

func setContentType(contentType string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", contentType)
		ctx.Next()
	}
}

func tracing() func(ctx *gin.Context) {
	textMapPropagator := otel.GetTextMapPropagator()

	return func(ctx *gin.Context) {
		ctx.Set(ginTracerKey, tracer.Span)

		savedCtx := ctx.Request.Context()
		defer func() {
			ctx.Request = ctx.Request.WithContext(savedCtx)
		}()

		// create ctx with propagators.
		ctxTracer := textMapPropagator.Extract(savedCtx, propagation.HeaderCarrier(ctx.Request.Header))

		path := ctx.FullPath()
		method := ctx.Request.Method
		if path == "" {
			path = fmt.Sprintf("HTTP %s route not found", method)
		}
		spanName := fmt.Sprintf("%s %s", method, path)

		// Create a span
		ctxTracer, span := tracer.Span(ctxTracer, spanName,
			tracer.SpanStartOption{Key: string(semconv.HTTPSchemeKey), Value: ctx.Request.URL.Scheme},
			tracer.SpanStartOption{Key: string(semconv.HTTPMethodKey), Value: method},
			tracer.SpanStartOption{Key: string(semconv.HTTPURLKey), Value: ctx.Request.URL.String()},
		)
		defer span.End()

		// pass the span through the request context
		ctx.Request = ctx.Request.WithContext(ctxTracer)

		// serve the request to the next middleware
		ctx.Next()

		status := ctx.Writer.Status()
		span.SetAttributes(semconv.HTTPStatusCode(status))

		if status >= 400 {
			span.SetStatus(codes.Error, "")
		}

		if len(ctx.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", ctx.Errors.String()))
		}

	}
}

func (r *ginRouter) Get(path string, f HandlerFunc) {
	r.router.GET(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Post(path string, f HandlerFunc) {
	r.router.POST(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Put(path string, f HandlerFunc) {
	r.router.PUT(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Delete(path string, f HandlerFunc) {
	r.router.DELETE(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Server(port string) error {
	return http.ListenAndServe(port, r.router)
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (m *ginRouter) ParseHandler(h http.HandlerFunc) HandlerFunc {
	return func(c Context) {
		h(c.GetResponseWriter(), c.GetRequestReader())
	}
}

type ginContext struct {
	r *gin.Context
	v validator.Validator
}

func newGinContext(ctx *gin.Context) Context {
	return &ginContext{
		r: ctx,
		v: validator.New(),
	}
}

func (c *ginContext) Context() context.Context {
	return c.r.Request.Context()
}

func (c *ginContext) JSON(statusCode int, data any) {
	c.r.JSON(statusCode, data)
}

func (c *ginContext) Decode(data any) error {
	return c.r.Bind(&data)
}

func (c *ginContext) GetQuery(query string) string {
	return c.r.Query(query)
}

func (c *ginContext) GetParam(param string) string {
	return c.r.Param(param)
}

func (c *ginContext) GetResponseWriter() http.ResponseWriter {
	return c.r.Writer
}

func (c *ginContext) GetRequestReader() *http.Request {
	return c.r.Request
}

func (c *ginContext) Validate(input any) error {
	if err := c.v.Validate(input); err != nil {
		return err.ToHttpErr()
	}
	return nil
}
