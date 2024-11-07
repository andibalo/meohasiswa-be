package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
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

type ThreadController struct {
	cfg     config.Config
	authSvc service.AuthService
}

func NewThreadController(cfg config.Config, authSvc service.AuthService) *ThreadController {

	return &ThreadController{
		cfg:     cfg,
		authSvc: authSvc,
	}
}

func (h *ThreadController) AddRoutes(r *gin.Engine) {
	ar := r.Group("/api/v1/thread")

	ar.POST("", h.CreateThread)
}

func (h *ThreadController) CreateThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.RegisterUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.Register(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Register] Failed to create user", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
