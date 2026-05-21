package facility

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/shopspring/decimal"
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
	HourlyRate  *float64 `json:"hourly_rate"`
	OwnerEmail  *string  `json:"owner_email"`
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

	ownerID := userID

	// Admin can create facilities for other owners
	userRole, _ := c.Get("role")
	if userRole != nil && userRole.(string) == "ADMIN" && dto.OwnerEmail != nil && *dto.OwnerEmail != "" {
		ownerUser, err := h.service.LookupUserByEmail(*dto.OwnerEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Owner not found with that email"})
			return
		}
		ownerID = ownerUser.ID
	}

	facility := &models.Facility{
		Name:        dto.Name,
		Description: dto.Description,
		Type:        dto.Type,
		Address:     dto.Address,
		City:        dto.City,
		Lat:         dto.Lat,
		Lng:         dto.Lng,
		OwnerID:     ownerID,
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

// GetMine returns facilities owned by the authenticated user
func (h *Handler) GetMine(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	facilities, err := h.service.GetByOwnerID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch facilities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": facilities})
}

// UpdateFacility updates a facility's properties (owner only)
func (h *Handler) UpdateFacility(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	facilityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}

	var body struct {
		Name               *string  `json:"name"`
		Description        *string  `json:"description"`
		Address            *string  `json:"address"`
		City               *string  `json:"city"`
		Type               *string  `json:"type"`
		HourlyRate         *float64 `json:"hourly_rate"`
		BookingMode        *string  `json:"booking_mode"`
		RequiresPrepayment *bool    `json:"requires_prepayment"`
		PrepaymentCost     *float64 `json:"prepayment_cost"`
		BankAccount        *string  `json:"bank_account"`
		TransferTitle      *string  `json:"transfer_title"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Description != nil {
		updates["description"] = body.Description
	}
	if body.Address != nil {
		updates["address"] = *body.Address
	}
	if body.City != nil {
		updates["city"] = *body.City
	}
	if body.Type != nil {
		updates["type"] = *body.Type
	}
	if body.HourlyRate != nil {
		updates["hourly_rate"] = *body.HourlyRate
	}
	if body.BookingMode != nil {
		validModes := map[string]bool{"INDIVIDUAL": true, "TEAM": true, "BOTH": true}
		if !validModes[*body.BookingMode] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "booking_mode must be INDIVIDUAL, TEAM, or BOTH"})
			return
		}
		updates["booking_mode"] = *body.BookingMode
	}
	if body.RequiresPrepayment != nil {
		updates["requires_prepayment"] = *body.RequiresPrepayment
	}
	if body.PrepaymentCost != nil {
		updates["prepayment_cost"] = decimal.NewFromFloat(*body.PrepaymentCost)
	}
	if body.BankAccount != nil {
		updates["bank_account"] = *body.BankAccount
	}
	if body.TransferTitle != nil {
		updates["transfer_title"] = *body.TransferTitle
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	userRole, _ := c.Get("role")
	isAdmin := userRole != nil && userRole.(string) == "ADMIN"

	if isAdmin {
		if err := h.service.UpdateByAdmin(facilityID, updates); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
			return
		}
	} else {
		if err := h.service.UpdateByOwner(facilityID, userID, updates); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found or not owned by you"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Facility updated"})
}

// UpdateSlots updates opening hours for a facility (owner only)
func (h *Handler) UpdateSlots(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	facilityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role != nil && role.(string) == "ADMIN"

	if !isAdmin {
		// Verify ownership
		facilities, err := h.service.GetByOwnerID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
			return
		}
		owned := false
		for _, f := range facilities {
			if f.ID == facilityID {
				owned = true
				break
			}
		}
		if !owned {
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found or not owned by you"})
			return
		}
	}

	var body struct {
		Slots []models.FacilitySlot `json:"slots" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSlots(facilityID, body.Slots); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slots updated"})
}

// ToggleClose closes or reopens a facility (owner only)
func (h *Handler) ToggleClose(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	facilityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}

	var body struct {
		ClosedUntil *string `json:"closed_until"` // ISO date string, null to reopen
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var closedUntil *time.Time
	if body.ClosedUntil != nil && *body.ClosedUntil != "" {
		t, err := time.Parse("2006-01-02", *body.ClosedUntil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
		closedUntil = &t
	}

	if err := h.service.SetClosed(facilityID, userID, closedUntil); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found or not owned by you"})
		return
	}

	if closedUntil == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Facility reopened"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Facility closed"})
	}
}
