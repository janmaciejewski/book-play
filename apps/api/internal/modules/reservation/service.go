package reservation

import (
	"errors"
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

// GetByFacilityOwnerID returns all reservations for facilities owned by the given user
func (s *Service) GetByFacilityOwnerID(ownerID uuid.UUID) ([]models.Reservation, error) {
	var facilityIDs []uuid.UUID
	s.db.Model(&models.Facility{}).Where("owner_id = ?", ownerID).Pluck("id", &facilityIDs)

	if len(facilityIDs) == 0 {
		return []models.Reservation{}, nil
	}

	var reservations []models.Reservation
	if err := s.db.
		Preload("Facility").
		Preload("User").
		Preload("Team").
		Where("facility_id IN ?", facilityIDs).
		Order("date DESC, start_time DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// UpdateReservationStatus updates status with ownership verification
func (s *Service) UpdateReservationStatus(reservationID, ownerID uuid.UUID, status models.ReservationStatus) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var reservation models.Reservation
		if err := tx.First(&reservation, "id = ?", reservationID).Error; err != nil {
			return err
		}

		// Verify the owner owns this facility
		var facility models.Facility
		if err := tx.First(&facility, "id = ? AND owner_id = ?", reservation.FacilityID, ownerID).Error; err != nil {
			return errors.New("you can only manage reservations for your own facilities")
		}

		if reservation.Status != models.StatusPending {
			return errors.New("only pending reservations can be updated")
		}

		return tx.Model(&reservation).Update("status", status).Error
	})
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

// IsUserTeamCaptain checks if the user is a captain of the given team
func (s *Service) IsUserTeamCaptain(userID, teamID uuid.UUID) (bool, error) {
	var count int64
	if err := s.db.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id = ? AND role = ?", teamID, userID, models.TeamRoleCaptain).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func parseDate(dateStr string) time.Time {
	// Simple date parsing - adjust format as needed
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
