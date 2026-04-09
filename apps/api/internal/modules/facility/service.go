package facility

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll(filterType, filterCity string) ([]models.Facility, error) {
	var facilities []models.Facility
	query := s.db.Model(&models.Facility{})

	if filterType != "" {
		query = query.Where("type = ?", filterType)
	}
	if filterCity != "" {
		query = query.Where("city ILIKE ?", "%"+filterCity+"%")
	}

	if err := query.Find(&facilities).Error; err != nil {
		return nil, err
	}
	return facilities, nil
}

func (s *Service) GetByID(id uuid.UUID) (*models.Facility, error) {
	var facility models.Facility
	if err := s.db.Preload("Owner").Preload("Slots").First(&facility, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}

func (s *Service) Create(facility *models.Facility) error {
	return s.db.Create(facility).Error
}

// GetAvailability returns available time slots for a facility on a given date
func (s *Service) GetAvailability(id uuid.UUID, date time.Time) ([]string, error) {
	// Get facility to check operating hours
	var facility models.Facility
	if err := s.db.First(&facility, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Get existing reservations for this date
	var reservations []models.Reservation
	if err := s.db.Where("facility_id = ? AND date = ? AND status != ?", id, date, models.StatusCancelled).
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	// Generate all possible hours (8:00 - 22:00)
	allHours := []string{}
	for h := 8; h < 22; h++ {
		allHours = append(allHours, fmt.Sprintf("%02d:00", h))
	}

	// Mark booked hours
	bookedHours := make(map[string]bool)
	for _, r := range reservations {
		startHour := r.StartTime[:5] // Get HH:MM
		bookedHours[startHour] = true
	}

	// Return available hours
	available := []string{}
	for _, hour := range allHours {
		if !bookedHours[hour] {
			available = append(available, hour)
		}
	}

	return available, nil
}
