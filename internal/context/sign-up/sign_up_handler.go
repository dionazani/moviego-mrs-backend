package contextsignup

import (
	"net/http"
	"time"

	dto "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SignUpHandler handles HTTP requests for user sign-up.
type SignUpHandler struct {
	signUpService SignUpService
}

// NewSignUpHandler creates a new instance of SignUpHandler.
func NewSignUpHandler(signUpService SignUpService) *SignUpHandler {
	return &SignUpHandler{
		signUpService: signUpService,
	}
}

// SignUp binds the JSON payload, calls the AddNew service function, and returns the response.
func (h *SignUpHandler) SignUp(c *gin.Context) {
	var req SignUpDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Invalid request payload: " + err.Error(),
			Data:            nil,
		})
		return
	}

	res, err := h.signUpService.AddNew(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Failed to sign up: " + err.Error(),
			Data:            nil,
		})
		return
	}

	c.JSON(res.ResponseStatus, res)
}

// LoadById handles the HTTP GET request to retrieve a user by ID.
func (h *SignUpHandler) LoadById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Invalid UUID format: " + err.Error(),
			Data:            nil,
		})
		return
	}

	res, err := h.signUpService.LoadById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Failed to retrieve user: " + err.Error(),
			Data:            nil,
		})
		return
	}

	c.JSON(res.ResponseStatus, res)
}
