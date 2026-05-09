package team

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"gorm.io/gorm"
)

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

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
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			return err
		}
		member := &models.TeamMember{
			TeamID: team.ID, UserID: team.CaptainID, Role: models.TeamRoleCaptain,
		}
		return tx.Create(member).Error
	})
}

func (s *Service) AddMember(teamID, userID uuid.UUID, role models.TeamRole) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var existing models.TeamMember
		if err := tx.Where("team_id = ? AND user_id = ?", teamID, userID).First(&existing).Error; err == nil {
			return errors.New("user is already a member of this team")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		member := &models.TeamMember{TeamID: teamID, UserID: userID, Role: role}
		return tx.Create(member).Error
	})
}

func (s *Service) RemoveMember(teamID, memberID, requestorID uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var member models.TeamMember
		if err := tx.First(&member, "id = ? AND team_id = ?", memberID, teamID).Error; err != nil {
			return err
		}
		if member.Role == models.TeamRoleCaptain {
			return errors.New("the captain cannot be removed; transfer the captain role first")
		}
		// Self-removal (leaving) is always allowed for non-captain
		if member.UserID == requestorID {
			return tx.Delete(&member).Error
		}
		return tx.Delete(&member).Error
	})
}

func (s *Service) UpdateMemberRole(teamID, memberID uuid.UUID, newRole models.TeamRole) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var member models.TeamMember
		if err := tx.First(&member, "id = ? AND team_id = ?", memberID, teamID).Error; err != nil {
			return err
		}
		if newRole == models.TeamRoleCaptain && member.Role != models.TeamRoleCaptain {
			tx.Model(&models.TeamMember{}).Where("team_id = ? AND role = ?", teamID, models.TeamRoleCaptain).Update("role", models.TeamRoleAdmin)
			tx.Model(&models.Team{}).Where("id = ?", teamID).Update("captain_id", member.UserID)
		}
		return tx.Model(&member).Update("role", newRole).Error
	})
}

func (s *Service) UpdateTeam(teamID uuid.UUID, dto *UpdateDTO) (*models.Team, error) {
	var team models.Team
	if err := s.db.First(&team, "id = ?", teamID).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if dto.Name != "" {
		updates["name"] = dto.Name
	}
	if dto.Description != nil {
		updates["description"] = dto.Description
	}
	if len(updates) > 0 {
		if err := s.db.Model(&team).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(teamID)
}

func (s *Service) UpdateLogo(teamID uuid.UUID, logoPath string) error {
	return s.db.Model(&models.Team{}).Where("id = ?", teamID).Update("logo", logoPath).Error
}

func (s *Service) SaveLogo(teamID uuid.UUID, fileData []byte, extension string) (string, error) {
	uploadsDir := "uploads/teams"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}
	filename := fmt.Sprintf("%s_%d%s", teamID.String(), time.Now().Unix(), extension)
	filePath := filepath.Join(uploadsDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("failed to write logo file: %w", err)
	}
	return "/" + filePath, nil
}

func (s *Service) IsMemberSelf(teamID, memberID, userID uuid.UUID) (bool, error) {
	var member models.TeamMember
	err := s.db.Where("id = ? AND team_id = ? AND user_id = ?", memberID, teamID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) IsUserCaptainOrAdmin(teamID, userID uuid.UUID) (bool, error) {
	var member models.TeamMember
	err := s.db.Where("team_id = ? AND user_id = ? AND role IN ?",
		teamID, userID, []models.TeamRole{models.TeamRoleCaptain, models.TeamRoleAdmin}).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) SearchUsers(query string, excludeTeamID uuid.UUID) ([]models.User, error) {
	var existingIDs []uuid.UUID
	s.db.Model(&models.TeamMember{}).Where("team_id = ?", excludeTeamID).Pluck("user_id", &existingIDs)
	var users []models.User
	q := s.db.Where("email ILIKE ?", query+"%")
	if len(existingIDs) > 0 {
		q = q.Where("id NOT IN ?", existingIDs)
	}
	if err := q.Limit(10).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetMyTeams(userID uuid.UUID) ([]models.Team, error) {
	var teamIDs []uuid.UUID
	s.db.Model(&models.TeamMember{}).Where("user_id = ?", userID).Pluck("team_id", &teamIDs)
	var teams []models.Team
	if len(teamIDs) > 0 {
		if err := s.db.Where("id IN ?", teamIDs).Preload("Members").Preload("Members.User").Find(&teams).Error; err != nil {
			return nil, err
		}
	}
	return teams, nil
}

// --- Recruitment ---

func (s *Service) SetRecruitmentOpen(teamID uuid.UUID, open bool) error {
	return s.db.Model(&models.Team{}).Where("id = ?", teamID).Update("recruitment_open", open).Error
}

func (s *Service) ApplyRecruitment(teamID, userID uuid.UUID, message string) error {
	var team models.Team
	if err := s.db.First(&team, "id = ?", teamID).Error; err != nil {
		return err
	}
	if !team.RecruitmentOpen {
		return errors.New("recruitment is not open for this team")
	}
	var existingMember models.TeamMember
	if err := s.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&existingMember).Error; err == nil {
		return errors.New("you are already a member of this team")
	}
	var existingApp models.TeamRecruitmentApplication
	if err := s.db.Where("team_id = ? AND user_id = ? AND status = ?", teamID, userID, "PENDING").First(&existingApp).Error; err == nil {
		return errors.New("you already have a pending application")
	}
	app := &models.TeamRecruitmentApplication{TeamID: teamID, UserID: userID, Message: message, Status: "PENDING"}
	return s.db.Create(app).Error
}

func (s *Service) GetApplications(teamID uuid.UUID) ([]models.TeamRecruitmentApplication, error) {
	var apps []models.TeamRecruitmentApplication
	if err := s.db.Where("team_id = ?", teamID).Preload("User").Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) HandleApplication(teamID, appID uuid.UUID, status string) error {
	if status != "ACCEPTED" && status != "REJECTED" {
		return errors.New("status must be ACCEPTED or REJECTED")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		var app models.TeamRecruitmentApplication
		if err := tx.First(&app, "id = ? AND team_id = ?", appID, teamID).Error; err != nil {
			return err
		}
		if app.Status != "PENDING" {
			return errors.New("application has already been processed")
		}
		if err := tx.Model(&app).Update("status", status).Error; err != nil {
			return err
		}
		if status == "ACCEPTED" {
			member := &models.TeamMember{TeamID: teamID, UserID: app.UserID, Role: models.TeamRoleMember}
			return tx.Create(member).Error
		}
		return nil
	})
}
