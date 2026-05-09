package user

import (
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"gorm.io/gorm"
)

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

func (s *Service) GetAll() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateRole(id uuid.UUID, role string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.User{}, "id = ?", id).Error
}

type UpdateProfileDTO struct {
	Phone             *string `json:"phone"`
	Bio               *string `json:"bio"`
	City              *string `json:"city"`
	Country           *string `json:"country"`
	Position          *string `json:"position"`
	PreferredPosition *string `json:"preferred_position"`
	Age               *int    `json:"age"`
	Avatar            *string `json:"avatar"`
}

func (s *Service) UpdateProfile(userID uuid.UUID, dto *UpdateProfileDTO) (*models.User, error) {
	updates := map[string]interface{}{}
	if dto.Phone != nil {
		updates["phone"] = dto.Phone
	}
	if dto.Bio != nil {
		updates["bio"] = dto.Bio
	}
	if dto.City != nil {
		updates["city"] = dto.City
	}
	if dto.Country != nil {
		updates["country"] = dto.Country
	}
	if dto.Position != nil {
		updates["position"] = dto.Position
	}
	if dto.PreferredPosition != nil {
		updates["preferred_position"] = dto.PreferredPosition
	}
	if dto.Age != nil {
		updates["age"] = dto.Age
	}
	if dto.Avatar != nil {
		updates["avatar"] = dto.Avatar
	}
	if len(updates) > 0 {
		if err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(userID)
}
