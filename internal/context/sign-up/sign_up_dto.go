package contextsignup

// SignUpDTO represents the request payload for the sign-up API.
type SignUpDTO struct {
	Fullname             string `json:"fullname"`
	Gender               string `json:"gender"`
	Email                string `json:"email"`
	MobilePhone          string `json:"mobilePhone"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}
