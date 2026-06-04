package contextsignup

import (
	"errors"
	"net/http"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/infrastructure/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	var dto SignUpDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Invalid request payload: " + err.Error(),
			Data:            nil,
		})
		return
	}

	person, err := h.signUpService.AddNew(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Failed to sign up: " + err.Error(),
			Data:            nil,
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseCode:    http.StatusCreated,
		ResponseMessage: "User signed up successfully",
		Data: gin.H{
			"appUserId": person.ID,
		},
	})
}

// LoadById handles the HTTP GET request to retrieve a user by ID.
func (h *SignUpHandler) LoadById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Invalid UUID format: " + err.Error(),
			Data:            nil,
		})
		return
	}

	person, err := h.signUpService.LoadById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dtos.Response{
				Timestamp:       time.Now().Format(time.RFC3339),
				ResponseCode:    http.StatusNotFound,
				ResponseMessage: "User not found",
				Data:            nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Failed to retrieve user: " + err.Error(),
			Data:            nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Timestamp:       time.Now().Format(time.RFC3339),
		ResponseCode:    http.StatusOK,
		ResponseMessage: "User retrieved successfully",
		Data:            person,
	})
}
