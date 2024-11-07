package trace

import (
	"context"
	otelIn "github.com/andibalo/meowhasiswa-be/pkg/trace/otel"
	"sync"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type ITracer interface {
	Start(ctx context.Context, packageName, spanName string) (context.Context, func())
	BuildAsyncSpanContext(ctx context.Context, operationName, threadIDKey string) (context.Context, func())
}

var mutex sync.RWMutex
var itracer ITracer

type Config struct {
	ServiceName          string
	CollectorURL         string
	CollectorEnvironment string
	Insecure             bool
	SkipperPaths         []string
	FragmentRatio        float64
}

type Tracer struct {
	itracer          ITracer
	otelShutdownFunc func(context.Context) error
	otelGinOptions   []otelgin.Option
}

func Init(ctx context.Context, cfg Config) (*Tracer, error) {
	mutex.RLock()

	otelObj := otelIn.New(cfg.ServiceName, cfg.CollectorURL, cfg.FragmentRatio).SetInsecure(true)

	if len(cfg.SkipperPaths) > 0 {
		otelObj.SetSkipperPaths(cfg.SkipperPaths)
	}

	shutdownFunc, otelGinOpts, err := otelObj.Build(ctx)

	if err != nil {
		return nil, err
	}

	tr := &Tracer{
		itracer:          &OtelAPM{},
		otelShutdownFunc: shutdownFunc,
		otelGinOptions:   otelGinOpts,
	}

	mutex.RUnlock()

	// set for global use
	mutex.Lock()
	itracer = tr.itracer
	mutex.Unlock()

	return tr, nil
}

func (t *Tracer) SetGinMiddleware(router *gin.Engine, svcName string) {
	if t != nil && router != nil {
		router.Use(otelgin.Middleware(svcName, t.otelGinOptions...))
	}
}

func (t *Tracer) Close(ctx context.Context) error {
	if t.otelShutdownFunc != nil {
		return t.otelShutdownFunc(ctx)
	}
	return nil
}

func Start(ctx context.Context, packageName, spanName string) (context.Context, func()) {
	mutex.RLock()
	defer mutex.RUnlock()
	if itracer != nil {
		return itracer.Start(ctx, packageName, spanName)
	}

	return ctx, func() {}
}

func BuildAsyncSpanContext(ctx context.Context, operationName string, threadIDKey string) (context.Context, func()) {
	mutex.RLock()
	defer mutex.RUnlock()
	if itracer != nil {
		return itracer.BuildAsyncSpanContext(ctx, operationName, threadIDKey)
	}
	return ctx, func() {}
}

//func SetPostgresWrapper(dsn string) (*sqlx.DB, error) {
//	db, err := otelsqlx.Open("postgres", dsn,
//		otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return db, nil
//}

func getItracer() ITracer {
	mutex.RLock()
	defer mutex.RUnlock()
	return itracer
}
