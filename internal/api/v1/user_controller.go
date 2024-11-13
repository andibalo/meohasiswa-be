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

func (h *UserController) TestLog(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.TestLog", "controller")
	//defer endFunc()

	h.cfg.Logger().Info("test log from core service")

	c.JSON(http.StatusOK, nil)
}
