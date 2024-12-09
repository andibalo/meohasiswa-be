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

type UserController struct {
	cfg     config.Config
	userSvc service.UserService
}

func NewUserController(cfg config.Config, userSvc service.UserService) *UserController {

	return &UserController{
		cfg:     cfg,
		userSvc: userSvc,
	}
}

func (h *UserController) AddRoutes(r *gin.Engine) {
	ur := r.Group("/api/v1/user")

	ur.GET("/profile", middleware.JwtMiddleware(h.cfg), h.GetUserProfile)
	ur.GET("/device", middleware.JwtMiddleware(h.cfg), h.GetUserDevices)
	ur.POST("/device/:user_id", middleware.JwtMiddleware(h.cfg), h.CreateUserDevice)
	ur.PATCH("/ban/:user_id", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.BanUser)
	ur.PATCH("/unban/:user_id", middleware.JwtMiddleware(h.cfg), middleware.IsAdminMiddleware(h.cfg), h.UnBanUser)
	ur.GET("/test", h.TestLog)
}

func (h *UserController) GetUserProfile(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.GetUserProfile", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetUserProfileReq

	data.UserID = claims.ID
	data.UserEmail = claims.Email
	user, err := h.userSvc.GetUserProfile(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetUserProfile] Failed to get user profile", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, user, nil)
	return
}

func (h *UserController) CreateUserDevice(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.CreateUserDevice", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CreateUserDeviceReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateUserDevice] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserID = c.Param("user_id")
	data.UserEmail = claims.Email

	err := h.userSvc.CreateUserDevice(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateUserDevice] Failed to create user device data", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *UserController) GetUserDevices(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.GetUserDevices", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetUserDevicesReq

	data.NotificationToken = c.Query("notification_token")
	data.UserID = c.Query("user_id")
	data.UserEmail = claims.Email

	userDevices, err := h.userSvc.GetUserDevices(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetUserDevices] Failed to get user devices data", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, userDevices, nil)
	return
}

func (h *UserController) BanUser(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.BanUser", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.BanUserReq

	data.BanUserID = c.Param("user_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email

	err := h.userSvc.BanUser(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[BanUser] Failed to ban user", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *UserController) UnBanUser(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.UnBanUser", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UnBanUserReq

	data.UnBanUserID = c.Param("user_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email

	err := h.userSvc.UnBanUser(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UnBanUser] Failed to unban user", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *UserController) TestLog(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.TestLog", "controller")
	//defer endFunc()

	h.cfg.Logger().Info("test log from core service")

	c.JSON(http.StatusOK, nil)
}
