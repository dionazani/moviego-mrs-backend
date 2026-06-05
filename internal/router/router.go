package router

import (
	"net/http"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	contextsignup "github.com/dionazani/moviego-mrs-backend/internal/context/sign-up"
	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/repository"
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

	// Initialize Dependencies
	personRepo := infrastructurerepository.NewAppPersonRepository(db)
	userRepo := infrastructurerepository.NewAppUserRepository(db)
	signUpService := contextsignup.NewSignUpService(db, personRepo, userRepo, cfg.MasterUserRoleRegular)
	signUpHandler := contextsignup.NewSignUpHandler(signUpService)

	// API Route Group
	api := r.Group("/api")
	{
		api.POST("/sign-up/v1/", signUpHandler.SignUp)
		api.GET("/sign-up/v1/:id", signUpHandler.LoadById)
	}

	return r
}
