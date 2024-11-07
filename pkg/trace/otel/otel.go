package otel

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"strings"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

var (
	tracerProvider     *otlptrace.Exporter
	once               sync.Once
	reusablePropagator propagation.TextMapPropagator
)

var defaultSkipperPaths = map[string]struct{}{
	"/kobokan/": {},
	"/metrics":  {},
	"/":         {},
}

type Option struct {
}

type otelBuilder struct {
	serviceName   string
	serviceURL    string
	insecure      bool
	skipperPaths  map[string]struct{}
	fragmentRatio float64
}

type ITracer interface {
	SetInsecure(insecure bool) ITracer
	SetSkipperPaths(skipperPaths []string) ITracer
	Build(ctx context.Context) (shutdownFunc func(context.Context) error, otelGinOptions []otelgin.Option, err error)
}

// New initialized otel trace configuration
func New(serviceName, serviceURL string, fragmentRatio float64) *otelBuilder {
	return &otelBuilder{
		serviceName:   serviceName,
		serviceURL:    serviceURL,
		insecure:      true,
		skipperPaths:  defaultSkipperPaths,
		fragmentRatio: fragmentRatio,
	}
}

func (b *otelBuilder) SetInsecure(insecure bool) *otelBuilder {
	b.insecure = insecure
	return b
}

func (b *otelBuilder) SetSkipperPaths(skipperPaths []string) *otelBuilder {
	b.skipperPaths = make(map[string]struct{})
	for _, d := range skipperPaths {
		b.skipperPaths[d] = struct{}{}
	}
	return b
}

func (b *otelBuilder) Build(ctx context.Context) (shutdownFunc func(context.Context) error, otelGinOptions []otelgin.Option, err error) {
	resources, err := b.setupResource(ctx)
	if err != nil {
		log.Println("could not set resources", err)
	}

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if b.insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	// Export to jaeger using HTTP
	//headers := map[string]string{
	//	"content-type": "application/json",
	//}
	//
	//tracerProvider, err = otlptrace.New(
	//	context.Background(),
	//	otlptracehttp.NewClient(
	//		otlptracehttp.WithEndpoint(b.serviceURL),
	//		otlptracehttp.WithHeaders(headers),
	//		otlptracehttp.WithInsecure(),
	//	),
	//)

	tracerProvider, err = otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(b.serviceURL),
			otlptracegrpc.WithHeaders(nil),
		),
	)
	if err != nil {
		log.Println("could not initialize exporter", err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.TraceIDRatioBased(b.fragmentRatio)),
			trace.WithBatcher(tracerProvider), // Will process exporter as batch
			trace.WithResource(resources),
		),
	)

	otelGinOpts := []otelgin.Option{
		b.filterGinRequest(),
	}

	once.Do(func() {
		reusablePropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	})

	otel.SetTextMapPropagator(reusablePropagator)

	return tracerProvider.Shutdown, otelGinOpts, nil
}

func (b *otelBuilder) setupResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", b.serviceName),
			attribute.String("library.language", "go"),
		),
	)
}

func (b *otelBuilder) filterGinRequest() otelgin.Option {
	return otelgin.WithFilter(func(r *http.Request) bool {
		// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/instrumentation/github.com/gin-gonic/gin/otelgin/v0.37.0/instrumentation/github.com/gin-gonic/gin/otelgin/gintrace.go#L57-L64
		// should return false for excluded / rejected path and vice versa
		if strings.Contains(r.URL.Path, "/swagger/") {
			return false
		}

		_, rejected := b.skipperPaths[r.URL.Path]
		return !rejected
	})
}

func InjectTraceHeader(ctx context.Context, header http.Header) {
	if reusablePropagator == nil {
		return
	}

	reusablePropagator.Inject(ctx, propagation.HeaderCarrier(header))
}

// ReadTraceID is helper to read TraceID from context
func ReadTraceID(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "00000000000000000000000000000000", "0000000000000000"
	}
	span := otelTrace.SpanFromContext(ctx)
	return span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()
}
