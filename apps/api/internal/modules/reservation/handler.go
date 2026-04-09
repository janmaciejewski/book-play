package reservation

import (
	"log"
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
// @Summary Get all reservations
// @Description Get list of all reservations
// @Tags reservations
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /reservations [get]
func (h *Handler) GetAll(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse userID from string to uuid.UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is admin
	userRole, _ := c.Get("role")
	isAdmin := userRole == models.RoleAdmin || userRole == string(models.RoleAdmin)

	log.Printf("GetAll reservations - userID: %s, role: %v, isAdmin: %v", userID, userRole, isAdmin)

	var reservations []models.Reservation
	if isAdmin {
		// Admins see all reservations
		reservations, err = h.service.GetAll()
	} else {
		// Users see their own reservations
		reservations, err = h.service.GetByUserID(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservations})
}

// GetByID godoc
// @Summary Get reservation by ID
// @Description Get single reservation details
// @Tags reservations
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /reservations/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	reservation, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservation})
}

// Create godoc
// @Summary Create a new reservation
// @Description Create a new reservation
// @Tags reservations
// @Security BearerAuth
// @Param request body CreateDTO true "Reservation data"
// @Success 201 {object} map[string]interface{}
// @Router /reservations [post]
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

	reservation, err := h.service.CreateFromDTO(&dto, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": reservation})
}

// Cancel godoc
// @Summary Cancel a reservation
// @Description Cancel a reservation (only by owner)
// @Tags reservations
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [put]
func (h *Handler) Cancel(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	// Get reservation to check ownership
	reservation, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Check if user owns the reservation
	if reservation.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only cancel your own reservations"})
		return
	}

	// Check if already cancelled
	if reservation.Status == models.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation is already cancelled"})
		return
	}

	// Cancel the reservation
	if err := h.service.UpdateStatus(id, models.StatusCancelled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled successfully"})
}
