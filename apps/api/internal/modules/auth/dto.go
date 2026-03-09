package auth

// RegisterDTO represents the registration request body
type RegisterDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// LoginDTO represents the login request body
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenDTO represents the refresh token request body
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

// ForgotPasswordDTO represents the forgot password request body
type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordDTO represents the reset password request body
type ResetPasswordDTO struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
