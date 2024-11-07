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

type AuthController struct {
	cfg     config.Config
	authSvc service.AuthService
}

func NewAuthController(cfg config.Config, authSvc service.AuthService) *AuthController {

	return &AuthController{
		cfg:     cfg,
		authSvc: authSvc,
	}
}

func (h *AuthController) AddRoutes(r *gin.Engine) {
	ar := r.Group("/api/v1/auth")

	ar.POST("/register", h.Register)
	ar.POST("/login", h.Login)
}

func (h *AuthController) Register(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.RegisterUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.Register(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().Error("[Register] Failed to create user", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *AuthController) Login(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.LoginUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().Error("[Login] Failed to login", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, token, nil)
	return
}
