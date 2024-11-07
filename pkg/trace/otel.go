package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type OtelAPM struct{}

func (t *OtelAPM) Start(ctx context.Context, packageName, spanName string) (context.Context, func()) {
	ctx, span := otel.Tracer(packageName).Start(ctx, spanName)
	return ctx, func() {
		span.End()
	}
}

func (t *OtelAPM) BuildAsyncSpanContext(ctx context.Context, operationName, threadIDKey string) (context.Context, func()) {
	spanContext := trace.SpanFromContext(ctx)
	ctx = trace.ContextWithSpan(context.WithValue(context.Background(), threadIDKey, ctx.Value(threadIDKey)), spanContext)

	ctx, span := otel.Tracer("async process").Start(ctx, operationName)
	return ctx, func() {
		span.End()
	}

}
