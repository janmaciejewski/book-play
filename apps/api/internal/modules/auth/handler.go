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

// Register
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

// Login
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
	c.SetCookie("refresh_token", response.RefreshToken, int(config.AppConfigInstance.JWT.RefreshTokenExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, response)
}

// RefreshToken
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
	c.SetCookie("refresh_token", response.RefreshToken, int(config.AppConfigInstance.JWT.RefreshTokenExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, response)
}

// Logout
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken != "" {
		h.service.Logout(refreshToken)
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// SendOTP
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

// VerifyOTP
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

// ResetPassword sends an OTP and then resets the password after OTP verification
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

// GetMe
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
