package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vashafaris/nucleus/internal/infrastructure"
	"github.com/vashafaris/nucleus/internal/interfaces/http/handler"
)

type Router struct {
	engine *gin.Engine
	infra  *infrastructure.Manager
}

// New creates a new router
func New(infra *infrastructure.Manager) *Router {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Router{
		engine: engine,
		infra:  infra,
	}
}

// Setup configures all routes
func (r *Router) Setup() {
	// Health check endpoints
	healthHandler := handler.NewHealthHandler(r.infra)

	r.engine.GET("/health", healthHandler.Check())
	r.engine.GET("/health/live", healthHandler.Live())
	r.engine.GET("/health/ready", healthHandler.Ready())

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Product routes will be added in Step 3
		v1.GET("/products", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Products endpoint - coming soon"})
		})
	}

	// Metrics endpoint for Prometheus
	r.engine.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Metrics endpoint - coming soon"})
	})
}

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
