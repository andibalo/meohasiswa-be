package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/andibalo/meowhasiswa-be/pkg/integration/notifsvc"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type notificationService struct {
	cfg         config.Config
	notifClient notifsvc.INotifSvc
}

func NewNotificationService(cfg config.Config, notifClient notifsvc.INotifSvc) NotificationService {

	return &notificationService{
		cfg:         cfg,
		notifClient: notifClient,
	}
}

func (s *notificationService) SendPushNotification(ctx context.Context, req request.SendPushNotificationReq) error {
	//ctx, endFunc := trace.Start(ctx, "NotificationService.SendPushNotification", "service")
	//defer endFunc()

	sendNotifReq := notifsvc.SendPushNotificationReq{
		NotificationTokens: req.NotificationTokens,
		Title:              req.Title,
		Content:            req.Content,
	}

	_, err := s.notifClient.SendPushNotification(ctx, sendNotifReq)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[SendPushNotification] Notification Service Error", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to send push notification")
	}

	return nil
}
