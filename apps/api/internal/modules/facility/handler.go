package facility

import (
	"net/http"
	"time"

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
// @Summary Get all facilities
// @Description Get list of all facilities
// @Tags facilities
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /facilities [get]
func (h *Handler) GetAll(c *gin.Context) {
	filterType := c.Query("type")
	filterCity := c.Query("city")

	facilities, err := h.service.GetAll(filterType, filterCity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch facilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": facilities})
}

// GetByID godoc
// @Summary Get facility by ID
// @Description Get single facility details
// @Tags facilities
// @Security BearerAuth
// @Param id path string true "Facility ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /facilities/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}

	facility, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": facility})
}

// CreateDTO for creating a facility
type CreateDTO struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Type        string   `json:"type" binding:"required"`
	Address     string   `json:"address" binding:"required"`
	City        string   `json:"city" binding:"required"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	HourlyRate  float64  `json:"hourly_rate" binding:"required"`
}

// Create godoc
// @Summary Create a new facility
// @Description Create a new facility (facility owners only)
// @Tags facilities
// @Security BearerAuth
// @Param request body CreateDTO true "Facility data"
// @Success 201 {object} map[string]interface{}
// @Router /facilities [post]
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

	facility := &models.Facility{
		Name:        dto.Name,
		Description: dto.Description,
		Type:        dto.Type,
		Address:     dto.Address,
		City:        dto.City,
		Lat:         dto.Lat,
		Lng:         dto.Lng,
		OwnerID:     userID,
	}

	if err := h.service.Create(facility); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create facility"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": facility})
}

// GetAvailability godoc
// @Summary Get facility availability
// @Description Get available time slots for a facility on a given date
// @Tags facilities
// @Param id path string true "Facility ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /facilities/{id}/availability [get]
func (h *Handler) GetAvailability(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	available, err := h.service.GetAvailability(id, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": available})
}
