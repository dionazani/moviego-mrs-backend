package contextsignup

// SignUpDTO represents the request payload for the sign-up API.
type SignUpDTO struct {
	Fullname             string `json:"fullname" binding:"required"`
	Gender               string `json:"gender" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	MobilePhone          string `json:"mobilePhone" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}
