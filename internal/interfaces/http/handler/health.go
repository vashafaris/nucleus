package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vashafaris/nucleus/internal/infrastructure"
)

type HealthHandler struct {
	infra *infrastructure.Manager
}

func NewHealthHandler(infra *infrastructure.Manager) *HealthHandler {
	return &HealthHandler{
		infra: infra,
	}
}

type HealthResponse struct {
	Status   string                 `json:"status"`
	Services map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *HealthHandler) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		healthErrors := h.infra.Health()

		response := HealthResponse{
			Status:   "ok",
			Services: make(map[string]ServiceInfo),
		}

		// API is always healthy if we reach this point
		response.Services["api"] = ServiceInfo{
			Status: "healthy",
		}

		// Check PostgreSQL
		if err, exists := healthErrors["postgres"]; exists && err != nil {
			response.Status = "degraded"
			response.Services["postgres"] = ServiceInfo{
				Status:  "unhealthy",
				Message: err.Error(),
			}
		} else {
			response.Services["postgres"] = ServiceInfo{
				Status: "healthy",
			}
		}

		// Check Redis
		if err, exists := healthErrors["redis"]; exists && err != nil {
			response.Status = "degraded"
			response.Services["redis"] = ServiceInfo{
				Status:  "unhealthy",
				Message: err.Error(),
			}
		} else {
			response.Services["redis"] = ServiceInfo{
				Status: "healthy",
			}
		}

		statusCode := http.StatusOK
		if response.Status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

func (h *HealthHandler) Live() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	}
}

func (h *HealthHandler) Ready() gin.HandlerFunc {
	return func(c *gin.Context) {
		healthErrors := h.infra.Health()

		if len(healthErrors) > 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"errors": healthErrors,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}
