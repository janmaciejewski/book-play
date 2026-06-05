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

// GetByOwnerID – zwraca wszystkie obiekty danego właściciela
func (s *Service) GetByOwnerID(ownerID uuid.UUID) ([]models.Facility, error) {
	var facilities []models.Facility
	if err := s.db.Preload("Slots").Where("owner_id = ?", ownerID).Find(&facilities).Error; err != nil {
		return nil, err
	}
	return facilities, nil
}

// UpdateByOwner – aktualizuje obiekt po weryfikacji właściciela
func (s *Service) UpdateByOwner(facilityID, ownerID uuid.UUID, updates map[string]interface{}) error {
	result := s.db.Model(&models.Facility{}).
		Where("id = ? AND owner_id = ?", facilityID, ownerID).
		Updates(updates)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// UpdateByAdmin – aktualizuje dowolny obiekt bez sprawdzania własności
func (s *Service) UpdateByAdmin(facilityID uuid.UUID, updates map[string]interface{}) error {
	result := s.db.Model(&models.Facility{}).
		Where("id = ?", facilityID).
		Updates(updates)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// UpdateSlots – zastępuje wszystkie sloty godzinowe obiektu
func (s *Service) UpdateSlots(facilityID uuid.UUID, slots []models.FacilitySlot) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("facility_id = ?", facilityID).Delete(&models.FacilitySlot{}).Error; err != nil {
			return err
		}
		for i := range slots {
			slots[i].ID = uuid.Nil // wymusza nowe ID dla slotu
			slots[i].FacilityID = facilityID
			if err := tx.Create(&slots[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SetClosed – zamyka lub otwiera ponownie dany obiekt
func (s *Service) SetClosed(facilityID, ownerID uuid.UUID, closedUntil *time.Time) error {
	result := s.db.Model(&models.Facility{}).
		Where("id = ? AND owner_id = ?", facilityID, ownerID).
		Update("closed_until", closedUntil)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
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

// LookupUserByEmail – wyszukuje użytkownika po emailu (dla admina)
func (s *Service) LookupUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAvailability – zwraca dostępne godziny dla obiektu w podanym dniu
func (s *Service) GetAvailability(id uuid.UUID, date time.Time) ([]string, error) {
	// Pobiera obiekt, by sprawdzić godziny otwarcia i status zamknięcia
	var facility models.Facility
	if err := s.db.First(&facility, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Normalizuje datę do lokalnej strefy czasowej
	requestDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	// Sprawdza czy obiekt jest zamknięty w żądanym dniu
	if facility.ClosedUntil != nil {
		closedDate := time.Date(facility.ClosedUntil.Year(), facility.ClosedUntil.Month(), facility.ClosedUntil.Day(), 0, 0, 0, 0, facility.ClosedUntil.Location())
		if !requestDate.After(closedDate) {
			return []string{}, nil
		}
	}

	// Pobiera istniejące rezerwacje na ten dzień
	var reservations []models.Reservation
	if err := s.db.Where("facility_id = ? AND date = ? AND status != ?", id, date, models.StatusCancelled).
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	// Generuje wszystkie możliwe godziny (8:00 – 22:00)
	allHours := []string{}
	for h := 8; h < 22; h++ {
		allHours = append(allHours, fmt.Sprintf("%02d:00", h))
	}

	// Oznacza zajęte godziny
	bookedHours := make(map[string]bool)
	for _, r := range reservations {
		startHour := r.StartTime[:5] // Wyciąga HH:MM
		bookedHours[startHour] = true
	}

	// Sprawdza czy zapytanie dotyczy dzisiejszego dnia
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	isToday := requestDate.Equal(today)

	// Zwraca tylko wolne godziny
	available := []string{}
	for _, hour := range allHours {
		if bookedHours[hour] {
			continue
		}
		// Dla dzisiaj filtruje godziny, które już minęły
		if isToday {
			hourInt := 0
			fmt.Sscanf(hour, "%d:00", &hourInt)
			if hourInt <= now.Hour() {
				continue
			}
		}
		available = append(available, hour)
	}

	return available, nil
}
