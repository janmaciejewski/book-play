package reservation

import (
	"time"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := s.db.Preload("Facility").Preload("User").Preload("Team").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (s *Service) GetByID(id uuid.UUID) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := s.db.Preload("Facility").Preload("User").Preload("Team").First(&reservation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (s *Service) GetByUserID(userID uuid.UUID) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := s.db.Preload("Facility").Preload("Team").Where("user_id = ?", userID).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (s *Service) Create(reservation *models.Reservation) error {
	return s.db.Create(reservation).Error
}

func (s *Service) UpdateStatus(id uuid.UUID, status models.ReservationStatus) error {
	return s.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}

type CreateDTO struct {
	FacilityID string          `json:"facility_id" binding:"required"`
	TeamID     *string         `json:"team_id"`
	Date       string          `json:"date" binding:"required"`
	StartTime  string          `json:"start_time" binding:"required"`
	EndTime    string          `json:"end_time" binding:"required"`
	TotalPrice decimal.Decimal `json:"total_price" binding:"required"`
	Notes      *string         `json:"notes"`
}

func (s *Service) CreateFromDTO(dto *CreateDTO, userID uuid.UUID) (*models.Reservation, error) {
	facilityID, err := uuid.Parse(dto.FacilityID)
	if err != nil {
		return nil, err
	}

	reservation := &models.Reservation{
		FacilityID: facilityID,
		UserID:     userID,
		Date:       parseDate(dto.Date),
		StartTime:  dto.StartTime,
		EndTime:    dto.EndTime,
		Status:     models.StatusPending,
		TotalPrice: dto.TotalPrice,
		Notes:      dto.Notes,
	}

	if dto.TeamID != nil {
		teamID, err := uuid.Parse(*dto.TeamID)
		if err == nil {
			reservation.TeamID = &teamID
		}
	}

	if err := s.Create(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func parseDate(dateStr string) time.Time {
	// Simple date parsing - adjust format as needed
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
