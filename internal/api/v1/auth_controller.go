package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/gin-gonic/gin"
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

	ar.POST("/register")
}

func (h *AuthController) Register(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "AuthController.Register", "controller")
	//defer endFunc()

	var data request.RegisterUserReq

	if err := c.BindJSON(&data); err != nil {
		//httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	c.JSON(http.StatusOK, nil)
}
