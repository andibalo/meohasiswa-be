package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/gin-gonic/gin"
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

	ur.GET("/test", h.TestLog)

}

func (h *UserController) TestLog(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UserController.TestLog", "controller")
	//defer endFunc()

	h.cfg.Logger().Info("test log from core service")

	c.JSON(http.StatusOK, nil)
}
