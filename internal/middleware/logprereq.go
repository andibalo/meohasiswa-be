package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/andibalo/meowhasiswa-be/pkg/httpclient"
	"github.com/andibalo/meowhasiswa-be/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
)

func LogPreReq(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// client app id
		var clientID = ctx.Request.Header.Get(httpclient.XClientID)
		payload, _ := io.ReadAll(ctx.Request.Body)

		//traceID, spanID := observ.ReadTraceID(ctx.Request.Context())

		compactPayload := &bytes.Buffer{}
		err := json.Compact(compactPayload, payload)
		if err != nil {
			compactPayload = bytes.NewBuffer(payload)
		}
		// set client id
		ctx.Set("x-client-id", clientID)
		//ctx.Set("trace.id", traceID)
		ctx.Set("path", ctx.Request.URL.Path)
		ctx.Set("method", ctx.Request.Method)

		// payload for log
		logger.InfoWithContext(ctx, "Interceptor Log",
			zap.Any("payload", compactPayload),
			//zap.Any("trace.id", traceID),
			//zap.Any("span.id", spanID),
		)

		// payload for otel
		ctx.Set("payload", string(payload))

		// set req body again
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(payload))

		ctx.Next()
	}
}
