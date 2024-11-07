package trace

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func TracerLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		span := trace.SpanFromContext(c.Request.Context())

		// http header attributes
		span.SetAttributes(
			attribute.String("content_type", c.Request.Header.Get("Content-Type")),
			attribute.String("content_length", c.Request.Header.Get("Content-Length")),
		)

		// x-attributes
		span.SetAttributes(
			attribute.String("x-forwarded-for", c.Request.Header.Get("X-Forwarded-For")),
			attribute.String("x-forwarded-host", c.Request.Header.Get("X-Forwarded-Host")),
			attribute.String("x-client-id", c.Request.Header.Get("X-Client-ID")),
			attribute.String("x-client-version", c.Request.Header.Get("X-Client-Version")),
			attribute.String("x-request-id", c.Request.Header.Get("X-Request-ID")),
		)

		ctype := c.Request.Header.Get("Content-Type")

		attrs := []attribute.KeyValue{
			attribute.String("request.uri", c.Request.RequestURI),
		}

		if strings.HasPrefix(ctype, "application/json") {
			// read request body content
			requestBody, err := io.ReadAll(c.Request.Body)

			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)

				c.Next()
				return
			}

			attrs = append(attrs, attribute.String("request.body", string(requestBody)))

			// restore request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		defer func() {
			span.AddEvent("Interceptor Log", trace.WithAttributes(attrs...))
		}()

		c.Next()

	}
}
