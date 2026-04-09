package user

import (
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
