package notifsvc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/pkg/httpclient"
	"github.com/andibalo/meowhasiswa-be/pkg/trace"
	"io"
)

type INotifSvc interface {
	CreateNotifTemplate(ctx context.Context, req CreateNotifTemplateReq) (res CreateNotifTemplateResp, err error)
}

type notifsvc struct {
	URL        string
	Token      string
	AppID      string
	httpClient httpclient.IHTTPClient
}

func NewNotificationService(cfg config.Config, httpClient httpclient.IHTTPClient) INotifSvc {
	return &notifsvc{
		URL:        cfg.GetNotifSvcCfg().URL,
		Token:      cfg.GetNotifSvcCfg().Token,
		AppID:      cfg.AppID(),
		httpClient: httpClient,
	}
}

func (ns *notifsvc) CreateNotifTemplate(ctx context.Context, req CreateNotifTemplateReq) (res CreateNotifTemplateResp, err error) {
	ctx, endFunc := trace.Start(ctx, "notifsvc.CreateNotifTemplate", "external")

	defer endFunc()

	resp, err := ns.httpClient.PostJSON(ctx, &httpclient.PropRequest{
		URI: ns.URL + "/api/v1/template",
		Headers: map[string]string{
			"X-App-Token": ns.Token,
			"X-Client-Id": ns.AppID,
		},
		Body: req,
	})

	if err != nil {
		return res, err
	}

	rawData, err := io.ReadAll(resp.Body)

	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if err = json.Unmarshal(rawData, &res); err != nil {
		return res, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 || resp.StatusCode == 402 {
			return res, errors.New("bad request")
		}

		if resp.StatusCode == 404 {
			return res, errors.New("not found")
		}

		return res, errors.New("internal server error")
	}

	return res, nil
}
