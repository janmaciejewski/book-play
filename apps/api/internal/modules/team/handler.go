package team

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAll godoc
// @Summary Get all teams
// @Description Get list of all teams
// @Tags teams
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /teams [get]
func (h *Handler) GetAll(c *gin.Context) {
	teams, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// GetByID godoc
// @Summary Get team by ID
// @Description Get single team details
// @Tags teams
// @Security BearerAuth
// @Param id path string true "Team ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /teams/{id} [get]
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

// CreateDTO for creating a team
type CreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// Create godoc
// @Summary Create a new team
// @Description Create a new team
// @Tags teams
// @Security BearerAuth
// @Param request body CreateDTO true "Team data"
// @Success 201 {object} map[string]interface{}
// @Router /teams [post]
func (h *Handler) Create(c *gin.Context) {
	userIDStr, _ := c.Get("userID")

	// Parse userID from string to uuid.UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team := &models.Team{
		Name:        dto.Name,
		Description: dto.Description,
		CaptainID:   userID,
	}

	if err := h.service.Create(team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	// Add creator as captain member
	h.service.AddMember(team.ID, userID, models.TeamRoleCaptain)

	c.JSON(http.StatusCreated, gin.H{"data": team})
}
