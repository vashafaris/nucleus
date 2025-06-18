package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status   string                 `json:"status"`
	Services map[string]ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status: "ok",
			Services: map[string]ServiceInfo{
				"api": {
					Status: "healthy",
				},
				// Will add more service checks in Step 2
			},
		}

		c.JSON(http.StatusOK, response)
	}
}
