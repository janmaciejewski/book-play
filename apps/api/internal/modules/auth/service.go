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

// Register tworzy nowe konto użytkownika z zahaszowanym hasłem
func (s *Service) Register(dto *RegisterDTO) (*models.User, error) {
	// Sprawdza czy użytkownik już istnieje w bazie
	var existingUser models.User
	if err := s.db.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hashuje hasło przed zapisaniem
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Zapisuje nowego użytkownika w bazie
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

// Login uwierzytelnia użytkownika i generuje tokeny JWT
func (s *Service) Login(dto *LoginDTO) (*AuthResponse, error) {
	// Szuka użytkownika po adresie email
	var user models.User
	if err := s.db.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Sprawdza czy konto jest aktywne
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Porównuje hasło z zapisanym hashem
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generuje access + refresh token
	return s.generateTokens(&user)
}

// RefreshToken wymienia refresh token na nową parę tokenów
func (s *Service) RefreshToken(tokenString string) (*AuthResponse, error) {
	// Wyszukuje ważny refresh token w bazie
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", tokenString, time.Now()).
		Preload("User").First(&refreshToken).Error; err != nil {
		return nil, ErrInvalidToken
	}

	// Sprawdza czy konto użytkownika jest nadal aktywne
	if !refreshToken.User.IsActive {
		return nil, ErrUserInactive
	}

	// Usuwa zużyty refresh token
	s.db.Delete(&refreshToken)

	// Tworzy nową parę tokenów
	return s.generateTokens(&refreshToken.User)
}

// Logout unieważnia refresh token użytkownika
func (s *Service) Logout(tokenString string) error {
	return s.db.Where("token = ?", tokenString).Delete(&models.RefreshToken{}).Error
}

// GetUserByID pobiera użytkownika po jego ID
func (s *Service) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// generateTokens tworzy access token (krótki) i refresh token (długi)
func (s *Service) generateTokens(user *models.User) (*AuthResponse, error) {
	cfg := config.AppConfigInstance.JWT

	// Tworzy krótkoterminowy access token JWT
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

	// Tworzy długoterminowy refresh token (UUID)
	refreshExpiry := time.Now().Add(cfg.RefreshTokenExpiry)
	refreshTokenString := uuid.New().String()

	// Zapisuje refresh token w bazie do późniejszej walidacji
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

// ResetPasswordWithOTP weryfikuje kod OTP i zmienia hasło użytkownika
func (s *Service) ResetPasswordWithOTP(email, code, newPassword string) error {
	var otp models.EmailVerificationOTP
	if err := s.db.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired verification code")
		}
		return fmt.Errorf("failed to verify OTP: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.db.Model(&models.User{}).Where("email = ?", email).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Usuwa wykorzystany kod OTP
	s.db.Delete(&otp)

	return nil
}

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
