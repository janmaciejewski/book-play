package team

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type CreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}
type UpdateDTO struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
type AddMemberDTO struct {
	Email string          `json:"email" binding:"required"`
	Role  models.TeamRole `json:"role" binding:"required"`
}
type UpdateMemberRoleDTO struct {
	Role models.TeamRole `json:"role" binding:"required"`
}

func (h *Handler) GetMyTeams(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	teams, err := h.service.GetMyTeams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

func (h *Handler) GetAll(c *gin.Context) {
	teams, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	team, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": team})
}

func (h *Handler) Create(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team := &models.Team{Name: dto.Name, Description: dto.Description, CaptainID: userID}
	if err := h.service.Create(team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": team})
}

func (h *Handler) Update(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the captain or team admin can update the team"})
		return
	}
	var dto UpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, err := h.service.UpdateTeam(teamID, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": team})
}

func (h *Handler) UploadLogo(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the captain or team admin can update the team logo"})
		return
	}
	file, header, err := c.Request.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No logo file provided"})
		return
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".svg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Allowed formats: jpg, jpeg, png, webp, svg"})
		return
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read logo file"})
		return
	}
	if len(fileData) > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file too large (max 5MB)"})
		return
	}
	logoPath, err := h.service.SaveLogo(teamID, fileData, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo"})
		return
	}
	if err := h.service.UpdateLogo(teamID, logoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team logo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"logo": logoPath}})
}

func (h *Handler) AddMember(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the captain or team admin can add members"})
		return
	}
	var dto AddMemberDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dto.Role != models.TeamRoleCaptain && dto.Role != models.TeamRoleAdmin && dto.Role != models.TeamRoleMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	users, err := h.service.SearchUsers(dto.Email, teamID)
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := h.service.AddMember(teamID, users[0].ID, dto.Role); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	team, _ := h.service.GetByID(teamID)
	c.JSON(http.StatusOK, gin.H{"data": team})
}

func (h *Handler) RemoveMember(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}
	isSelf, _ := h.service.IsMemberSelf(teamID, memberID, userID)
	if !isSelf {
		if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the captain or team admin can remove members"})
			return
		}
	}
	if err := h.service.RemoveMember(teamID, memberID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, _ := h.service.GetByID(teamID)
	c.JSON(http.StatusOK, gin.H{"data": team})
}

func (h *Handler) UpdateMemberRole(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the captain or team admin can change member roles"})
		return
	}
	var dto UpdateMemberRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dto.Role != models.TeamRoleCaptain && dto.Role != models.TeamRoleAdmin && dto.Role != models.TeamRoleMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	if err := h.service.UpdateMemberRole(teamID, memberID, dto.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	team, _ := h.service.GetByID(teamID)
	c.JSON(http.StatusOK, gin.H{"data": team})
}

// --- Recruitment ---
func (h *Handler) ToggleRecruitment(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, _ := uuid.Parse(c.Param("id"))
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only captain/admin can manage recruitment"})
		return
	}
	var body struct {
		Open bool `json:"open"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.service.SetRecruitmentOpen(teamID, body.Open)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"recruitment_open": body.Open}})
}

func (h *Handler) ApplyRecruitment(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, _ := uuid.Parse(c.Param("id"))
	var body struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ApplyRecruitment(teamID, userID, body.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Application submitted"})
}

func (h *Handler) GetApplications(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, _ := uuid.Parse(c.Param("id"))
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only captain/admin can view applications"})
		return
	}
	apps, _ := h.service.GetApplications(teamID)
	c.JSON(http.StatusOK, gin.H{"data": apps})
}

func (h *Handler) HandleApplication(c *gin.Context) {
	userID, _ := h.getUserID(c)
	teamID, _ := uuid.Parse(c.Param("id"))
	appID, _ := uuid.Parse(c.Param("appId"))
	if ok, _ := h.service.IsUserCaptainOrAdmin(teamID, userID); !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only captain/admin can handle applications"})
		return
	}
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.HandleApplication(teamID, appID, body.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Application " + body.Status})
}

func (h *Handler) SearchUsers(c *gin.Context) {
	teamID, _ := uuid.Parse(c.Param("id"))
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query 'q' required"})
		return
	}
	users, _ := h.service.SearchUsers(query, teamID)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, nil
	}
	return uuid.Parse(userIDStr.(string))
}
