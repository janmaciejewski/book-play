package team

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

func (s *Service) GetAll() ([]models.Team, error) {
	var teams []models.Team
	if err := s.db.Preload("Members").Preload("Members.User").Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *Service) GetByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	if err := s.db.Preload("Members").Preload("Members.User").First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *Service) Create(team *models.Team) error {
	return s.db.Create(team).Error
}

func (s *Service) AddMember(teamID, userID uuid.UUID, role models.TeamRole) error {
	member := &models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}
	return s.db.Create(member).Error
}
