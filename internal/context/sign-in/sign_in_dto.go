package signin

// SignInRequest defines the payload expected for a sign-in operation.
type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenData holds the JWT tokens returned upon successful sign-in.
type TokenData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh-token"`
}
