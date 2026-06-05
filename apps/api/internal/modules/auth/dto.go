package auth

// RegisterDTO – dane wejściowe dla rejestracji użytkownika
type RegisterDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// LoginDTO – dane wejściowe dla logowania
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenDTO – dane do odświeżenia pary tokenów
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse – odpowiedź z tokenami po uwierzytelnieniu
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

// ForgotPasswordDTO – dane do wysłania kodu resetu hasła
type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordDTO – dane do zmiany hasła
type ResetPasswordDTO struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
