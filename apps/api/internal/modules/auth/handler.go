package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/mail"
)

type Handler struct {
	service     *Service
	mailService *mail.Service
}

func NewHandler(service *Service, mailService *mail.Service) *Handler {
	return &Handler{service: service, mailService: mailService}
}

// Register – rejestracja nowego użytkownika
func (h *Handler) Register(c *gin.Context) {
	var dto RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.Register(&dto)
	if err != nil {
		if err == ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": user.ID, "email": user.Email, "first_name": user.FirstName, "last_name": user.LastName, "role": user.Role,
	})
}

// Login – logowanie i zwrot tokenów JWT
func (h *Handler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Login(&dto)
	if err != nil {
		if err == ErrInvalidCredentials || err == ErrUserInactive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/",
		MaxAge:   int(config.AppConfigInstance.JWT.RefreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, response)
}

// RefreshToken – odświeżenie pary tokenów
func (h *Handler) RefreshToken(c *gin.Context) {
	var dto RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil || dto.RefreshToken == "" {
		dto.RefreshToken, _ = c.Cookie("refresh_token")
		if dto.RefreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}
	}
	response, err := h.service.RefreshToken(dto.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/",
		MaxAge:   int(config.AppConfigInstance.JWT.RefreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, response)
}

// Logout – wylogowanie i unieważnienie refresh tokena
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken != "" {
		h.service.Logout(refreshToken)
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// SendOTP – wysyła kod weryfikacyjny na email
func (h *Handler) SendOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.mailService.GenerateAndSendOTP(body.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

// VerifyOTP – sprawdza poprawność kodu weryfikacyjnego
func (h *Handler) VerifyOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	valid, err := h.mailService.VerifyOTP(body.Email, body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify code"})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}

// ResetPassword – zmienia hasło po weryfikacji kodem OTP
func (h *Handler) ResetPassword(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required,len=6"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ResetPasswordWithOTP(body.Email, body.Code, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hasło zostało zmienione"})
}

// GetMe – zwraca dane zalogowanego użytkownika
func (h *Handler) GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": user.ID, "email": user.Email, "first_name": user.FirstName, "last_name": user.LastName, "role": user.Role,
	})
}
