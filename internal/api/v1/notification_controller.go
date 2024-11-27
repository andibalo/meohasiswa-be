package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/middleware"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type NotificationController struct {
	cfg      config.Config
	notifSvc service.NotificationService
}

func NewNotificationController(cfg config.Config, notifSvc service.NotificationService) *NotificationController {

	return &NotificationController{
		cfg:      cfg,
		notifSvc: notifSvc,
	}
}

func (h *NotificationController) AddRoutes(r *gin.Engine) {
	nr := r.Group("/api/v1/notification")

	nr.POST("/push", middleware.JwtMiddleware(h.cfg), h.SendPushNotification)
}

func (h *NotificationController) SendPushNotification(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "NotificationController.SendPushNotification", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.SendPushNotificationReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[SendPushNotification] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.notifSvc.SendPushNotification(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[SendPushNotification] Failed to send push notification", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrBadRequest))
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)

	return
}
