package httpclient

import (
	"context"
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"time"

	resty "github.com/go-resty/resty/v2"
)

const (
	defaultTimeout  = 60 * time.Second
	infoRequestDump = `httpclient sent request: client_uri=%s request=%s`

	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "github.com/andibalo/meowhasiswa-be/pkg/httpclient"
)

//go:generate mockery --name=IHTTPClient --case underscore
type IHTTPClient interface {
	Get(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Post(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Put(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Delete(ctx context.Context, prop *PropRequest) (*http.Response, error)

	GetJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PostJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PutJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PatchJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
}

type restyC struct {
	client *resty.Client
	config config.Config
}

type Options struct {
	Config config.Config
}

func Init(params Options) IHTTPClient {
	r := &restyC{
		client: resty.New(),
		config: params.Config,
	}

	r.initClient()

	return r
}

func (c *restyC) initClient() {
	rl := restyLogger{c.config.Logger()}

	timeout := defaultTimeout
	if c.config.HttpExternalServiceTimeout() > 0 {
		timeout = time.Duration(c.config.HttpExternalServiceTimeout()) * time.Second
	}

	c.client.SetTimeout(timeout).
		SetContentLength(true).
		SetCloseConnection(false).
		SetJSONEscapeHTML(true).
		SetDoNotParseResponse(true).
		SetHeader(`User-Agent`, "Meowhasiswa").
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			//append header x-request-id
			reqID := req.Context().Value(RequestID)
			if reqID != nil {
				req.SetHeader(RequestID, reqID.(string))
			}

			if len(req.Header[XClientID]) < 1 {
				req.SetHeader(XClientID, c.config.AppName())
			}

			if len(req.Header[XClientVersion]) < 1 {
				req.SetHeader(XClientVersion, c.config.AppVersion())
			}

			return nil
		}).
		SetPreRequestHook(func(client *resty.Client, req *http.Request) error {
			//dump request
			if req != nil {
				uri := req.URL

				dumpRequest, err := httputil.DumpRequestOut(req, true)
				if err != nil {
					c.config.Logger().Error("error dump request", zap.Error(err))
				} else {
					c.config.Logger().Info(fmt.Sprintf(infoRequestDump, uri, string(dumpRequest)))
				}
			}

			return nil
		}).
		SetLogger(rl)
}
