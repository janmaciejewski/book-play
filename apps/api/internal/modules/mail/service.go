package mail

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"gorm.io/gorm"
)

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

func (s *Service) GenerateAndSendOTP(email string) error {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	otp := &models.EmailVerificationOTP{
		Email: email, Code: code, ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.db.Create(otp).Error; err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}
	return s.sendOTPEmail(email, code)
}

func (s *Service) VerifyOTP(email, code string) (bool, error) {
	var otp models.EmailVerificationOTP
	err := s.db.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).First(&otp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	s.db.Delete(&otp)
	return true, nil
}

func (s *Service) sendOTPEmail(to, code string) error {
	cfg := config.AppConfigInstance.SMTP
	subject := "BookPlay - Kod weryfikacyjny"
	body := fmt.Sprintf("Witaj!\r\n\r\nTwój kod weryfikacyjny dla BookPlay to: %s\r\n\r\nKod jest ważny przez 10 minut.\r\n\r\nPozdrawiamy,\r\nZespół BookPlay", code)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", cfg.From, to, subject, body))
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := smtp.SendMail(addr, auth, cfg.From, []string{to}, msg); err != nil {
		return fmt.Errorf("SMTP error: %w", err)
	}
	return nil
}
