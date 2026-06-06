package router

import (
	"net/http"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	contextsignin "github.com/dionazani/moviego-mrs-backend/internal/context/sign-in"
	contextsignup "github.com/dionazani/moviego-mrs-backend/internal/context/sign-up"
	infrastructurerepository "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/repository"
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

	// Initialize Repositories
	personRepo := infrastructurerepository.NewAppPersonRepository(db)
	userRepo := infrastructurerepository.NewAppUserRepository(db)
	userTokenRepo := infrastructurerepository.NewAppUserTokenRepository(db)

	// Initialize Services
	signUpService := contextsignup.NewSignUpService(db, personRepo, userRepo, cfg.MasterUserRoleRegular)
	signInService := contextsignin.NewSignInService(personRepo, userRepo, userTokenRepo, cfg.JwtSecret, cfg.LoginLockoutDuration)

	// Initialize Handlers
	signUpHandler := contextsignup.NewSignUpHandler(signUpService)
	signInHandler := contextsignin.NewSignInHandler(signInService)

	// API Route Group
	api := r.Group("/api")
	{
		// Sign-Up endpoints
		api.POST("/sign-up/v1/", signUpHandler.SignUp)
		api.GET("/sign-up/v1/:id", signUpHandler.LoadById)

		// Sign-In endpoint
		api.POST("/sign-in/v1", signInHandler.SignIn)
	}

	return r
}
