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
	ar.POST("/verify-email", h.VerifyEmail)
	ar.PATCH("/reset-password", h.ResetPassword)
	ar.PATCH("/reset-password/code/verify", h.VerifyResetPassword)
}

func (h *AuthController) Register(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.RegisterUserReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Register] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.Register(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Register] Failed to register user", zap.Error(err))
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
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Login] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	token, err := h.authSvc.Login(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[Login] Failed to login", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, token, nil)
	return
}

func (h *AuthController) VerifyEmail(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.VerifyEmail", "controller")
	//defer endFunc()

	var data request.VerifyEmailReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[VerifyEmail] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.VerifyEmail(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[VerifyEmail] Failed to verify email", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *AuthController) ResetPassword(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.ResetPassword", "controller")
	//defer endFunc()

	var data request.ResetPasswordReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ResetPassword] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.ResetPassword(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ResetPassword] Failed to reset password", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *AuthController) VerifyResetPassword(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.VerifyResetPassword", "controller")
	//defer endFunc()

	var data request.VerifyResetPasswordReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[VerifyResetPassword] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.authSvc.VerifyResetPassword(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[VerifyResetPassword] Failed to verify reset password code", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
