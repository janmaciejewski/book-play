package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUserInactive       = errors.New("user account is inactive")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Register creates a new user
func (s *Service) Register(dto *RegisterDTO) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        dto.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Role:         models.RolePlayer,
		IsActive:     true,
	}

	if dto.Phone != "" {
		user.Phone = &dto.Phone
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(dto *LoginDTO) (*AuthResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	return s.generateTokens(&user)
}

// RefreshToken generates new tokens from refresh token
func (s *Service) RefreshToken(tokenString string) (*AuthResponse, error) {
	// Find refresh token
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", tokenString, time.Now()).
		Preload("User").First(&refreshToken).Error; err != nil {
		return nil, ErrInvalidToken
	}

	// Check if user is active
	if !refreshToken.User.IsActive {
		return nil, ErrUserInactive
	}

	// Delete old refresh token
	s.db.Delete(&refreshToken)

	// Generate new tokens
	return s.generateTokens(&refreshToken.User)
}

// Logout invalidates the refresh token
func (s *Service) Logout(tokenString string) error {
	return s.db.Where("token = ?", tokenString).Delete(&models.RefreshToken{}).Error
}

// generateTokens creates access and refresh tokens
func (s *Service) generateTokens(user *models.User) (*AuthResponse, error) {
	cfg := config.AppConfigInstance.JWT

	// Generate access token
	accessExpiry := time.Now().Add(cfg.AccessTokenExpiry)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   accessExpiry.Unix(),
		"iat":   time.Now().Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshExpiry := time.Now().Add(cfg.RefreshTokenExpiry)
	refreshTokenString := uuid.New().String()

	// Store refresh token in database
	refreshToken := &models.RefreshToken{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: refreshExpiry,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(cfg.AccessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates an access token and returns the claims
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	cfg := config.AppConfigInstance.JWT

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
