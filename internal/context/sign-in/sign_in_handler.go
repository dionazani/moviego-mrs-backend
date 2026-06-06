package signin

import (
	"net/http"
	"time"

	dto "github.com/dionazani/moviego-mrs-backend/internal/infrastructure/dto"
	"github.com/gin-gonic/gin"
)

// SignInHandler handles HTTP requests for user sign-in.
type SignInHandler struct {
	signInService SignInService
}

// NewSignInHandler creates a new instance of SignInHandler.
func NewSignInHandler(signInService SignInService) *SignInHandler {
	return &SignInHandler{
		signInService: signInService,
	}
}

// SignIn binds the JSON payload, calls the SignIn service, and returns the response.
func (h *SignInHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseStatus:  http.StatusBadRequest,
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Invalid request payload: " + err.Error(),
			Data:            nil,
		})
		return
	}

	res, err := h.signInService.SignIn(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Timestamp:       time.Now().Format(time.RFC3339),
			ResponseStatus:  http.StatusInternalServerError,
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "Failed to sign in: " + err.Error(),
			Data:            nil,
		})
		return
	}

	c.JSON(res.ResponseStatus, res)
}
