package router

import (
	"net/http"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter initializes the Gin engine and registers all system endpoints
func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// Setup Gin router
	r := gin.Default()

	// Root Endpoint
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Movie Reservation System API!")
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"message": "MovieGo MRS Backend is fully operational",
		})
	})

	// Placeholder Route Groups for Contexts
	api := r.Group("/api")
	{
		// Authentication Context endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Register endpoint placeholder"})
			})
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Login endpoint placeholder"})
			})
			auth.POST("/activate", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Activation endpoint placeholder"})
			})
		}

		// Reservation Context endpoints
		reservations := api.Group("/reservations")
		{
			reservations.POST("/reserve", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Reserve seat endpoint placeholder"})
			})
		}
	}

	return r
}
