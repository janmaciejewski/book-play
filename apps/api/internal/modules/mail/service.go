package mail

import (
	"encoding/base64"
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

// encodeSubject koduje temat wiadomości w UTF-8 zgodnie z RFC 2047
func encodeSubject(s string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(s)) + "?="
}

func (s *Service) SendReservationReminder(reservation *models.Reservation) error {
	cfg := config.AppConfigInstance.SMTP

	// Formatuje datę w języku polskim
	months := []string{
		"stycznia", "lutego", "marca", "kwietnia", "maja", "czerwca",
		"lipca", "sierpnia", "wrze\xc5\x9bnia", "pa\xc5\xbadziernika", "listopada", "grudnia",
	}
	d := reservation.Date
	dateStr := fmt.Sprintf("%d %s %d", d.Day(), months[d.Month()-1], d.Year())

	facilityName := "obiekt"
	if reservation.Facility.Name != "" {
		facilityName = reservation.Facility.Name
	}

	facilityType := reservation.Facility.Type
	typeNames := map[string]string{
		"FOOTBALL":   "Pi\xc5\x82ka no\xc5\xbcna",
		"BASKETBALL": "Koszyk\xc3\xb3wka",
		"TENNIS":     "Tenis",
		"VOLLEYBALL": "Siatk\xc3\xb3wka",
		"SWIMMING":   "P\xc5\x82ywanie",
		"OTHER":      "Inne",
	}
	typeName, ok := typeNames[facilityType]
	if !ok {
		typeName = facilityType
	}

	subject := "BookPlay - Przypomnienie o rezerwacji"

	body := fmt.Sprintf(
		"Dzie\xc5\x84 dobry, %s!\r\n\r\n"+
			"Przypominamy o Twojej nadchodz\xc4\x85cej rezerwacji w BookPlay.\r\n\r\n"+
			"\xf0\x9f\x93\x8b Szczeg\xc3\xb3\xc5\x82y rezerwacji:\r\n"+
			"\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\r\n"+
			"\xf0\x9f\x94\xb9 Obiekt: %s\r\n"+
			"\xf0\x9f\x94\xb9 Typ: %s\r\n"+
			"\xf0\x9f\x94\xb9 Adres: %s, %s\r\n"+
			"\xf0\x9f\x94\xb9 Data: %s\r\n"+
			"\xf0\x9f\x94\xb9 Godzina: %s \xe2\x80\x93 %s\r\n"+
			"\xf0\x9f\x94\xb9 Cena: %s PLN\r\n",
		reservation.User.FirstName,
		facilityName,
		typeName,
		reservation.Facility.Address,
		reservation.Facility.City,
		dateStr,
		reservation.StartTime,
		reservation.EndTime,
		reservation.TotalPrice.String(),
	)

	if reservation.Notes != nil && *reservation.Notes != "" {
		body += fmt.Sprintf("\xf0\x9f\x94\xb9 Uwagi: %s\r\n", *reservation.Notes)
	}

	body += "\r\n" +
		"\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\xe2\x94\x81\r\n\r\n" +
		"Prosimy o punktualne przybycie na obiekt. W przypadku rezygnacji prosimy o " +
		"anulowanie rezerwacji w panelu BookPlay.\r\n\r\n" +
		"Dzi\xc4\x99kujemy za skorzystanie z BookPlay!\r\n\r\n" +
		"Z powa\xc5\xbcaniem,\r\n" +
		"Zesp\xc3\xb3\xc5\x82 BookPlay\r\n" +
		"\xf0\x9f\x93\xa7 bookplay.grodzisk@gmail.com"

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.From, reservation.User.Email, encodeSubject(subject), body,
	))

	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := smtp.SendMail(addr, auth, cfg.User, []string{reservation.User.Email}, msg); err != nil {
		return fmt.Errorf("SMTP error: %w", err)
	}

	// Oznacza przypomnienie jako wysłane
	if err := s.db.Model(reservation).Update("reminder_sent", true).Error; err != nil {
		return fmt.Errorf("failed to mark reminder sent: %w", err)
	}

	return nil
}

func (s *Service) ProcessReminders() error {
	// Wyszukuje rezerwacje za 3 dni, dla których jeszcze nie wysłano przypomnienia
	threeDaysFromNow := time.Now().AddDate(0, 0, 3)
	startOfDay := time.Date(threeDaysFromNow.Year(), threeDaysFromNow.Month(), threeDaysFromNow.Day(), 0, 0, 0, 0, threeDaysFromNow.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var reservations []models.Reservation
	if err := s.db.Preload("User").Preload("Facility").
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Where("status IN ?", []models.ReservationStatus{models.StatusConfirmed, models.StatusPending}).
		Where("reminder_sent = ?", false).
		Find(&reservations).Error; err != nil {
		return fmt.Errorf("failed to fetch reservations for reminders: %w", err)
	}

	for i := range reservations {
		if err := s.SendReservationReminder(&reservations[i]); err != nil {
			fmt.Printf("Failed to send reminder for reservation %s: %v\r\n", reservations[i].ID, err)
			continue
		}
		fmt.Printf("Sent reminder for reservation %s to %s\r\n", reservations[i].ID, reservations[i].User.Email)
	}

	return nil
}

func (s *Service) sendOTPEmail(to, code string) error {
	cfg := config.AppConfigInstance.SMTP
	fmt.Printf("[OTP EMAIL] Sending OTP to %s via SMTP host=%s port=%s user=%s\n", to, cfg.Host, cfg.Port, cfg.User)
	subject := "BookPlay - Kod weryfikacyjny"
	body := fmt.Sprintf("Witaj!\r\n\r\nTw\xc3\xb3j kod weryfikacyjny dla BookPlay to: %s\r\n\r\nKod jest wa\xc5\xbcny przez 10 minut.\r\n\r\nPozdrawiamy,\r\nZesp\xc3\xb3\xc5\x82 BookPlay", code)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", cfg.From, to, encodeSubject(subject), body))
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := smtp.SendMail(addr, auth, cfg.User, []string{to}, msg); err != nil {
		fmt.Printf("[OTP EMAIL] SMTP send failed: %v\n", err)
		return fmt.Errorf("SMTP error: %w", err)
	}
	fmt.Printf("[OTP EMAIL] Successfully sent OTP to %s\n", to)
	return nil
}
