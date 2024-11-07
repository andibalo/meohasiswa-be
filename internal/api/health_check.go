package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	healthCheckPath = "/health"
)

// HealthCheck is a standard, simple health check
type HealthCheck struct{}

// AddRoutes adds the routers for this API to the provided router (or subrouter)
func (h *HealthCheck) AddRoutes(r *gin.Engine) {
	r.GET(healthCheckPath, h.handler)
}

func (h *HealthCheck) handler(c *gin.Context) {

	c.String(http.StatusOK, "ok")
}
