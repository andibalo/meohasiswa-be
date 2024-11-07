package api

import "github.com/gin-gonic/gin"

type Handler interface {
	AddRoutes(r *gin.Engine)
}
