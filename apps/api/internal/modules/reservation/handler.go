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

func (h *Handler) GetAll(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Konwertuje ID użytkownika z stringa na UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Sprawdza czy użytkownik jest adminem
	userRole, _ := c.Get("role")
	isAdmin := userRole == models.RoleAdmin || userRole == string(models.RoleAdmin)

	log.Printf("GetAll reservations - userID: %s, role: %v, isAdmin: %v", userID, userRole, isAdmin)

	var reservations []models.Reservation
	if isAdmin {
		// Admin widzi wszystkie rezerwacje
		reservations, err = h.service.GetAll()
	} else {
		// Zwykły użytkownik widzi tylko swoje rezerwacje
		reservations, err = h.service.GetByUserID(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservations})
}

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

func (h *Handler) Create(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
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

	// Sprawdza tryb rezerwacji obiektu (indywidualny/drużynowy)
	facilityID, err := uuid.Parse(dto.FacilityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facility ID"})
		return
	}
	isTeamReservation := dto.TeamID != nil && *dto.TeamID != ""
	if err := h.service.CheckBookingMode(facilityID, isTeamReservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Jeśli podano team_id, weryfikuje czy użytkownik jest kapitanem drużyny
	if dto.TeamID != nil && *dto.TeamID != "" {
		teamID, err := uuid.Parse(*dto.TeamID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
			return
		}
		isCaptain, err := h.service.IsUserTeamCaptain(userID, teamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify team membership"})
			return
		}
		if !isCaptain {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the team captain can make reservations for the team"})
			return
		}
	}

	reservation, err := h.service.CreateFromDTO(&dto, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": reservation})
}

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

	// Pobiera rezerwację, aby sprawdzić własność
	reservation, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Sprawdza czy użytkownik jest właścicielem rezerwacji
	if reservation.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only cancel your own reservations"})
		return
	}

	// Sprawdza czy rezerwacja jest już anulowana
	if reservation.Status == models.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation is already cancelled"})
		return
	}

	// Anuluje rezerwację
	if err := h.service.UpdateStatus(id, models.StatusCancelled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled successfully"})
}

// GetForFacilityOwner – zwraca rezerwacje dla obiektów zalogowanego właściciela (lub wszystkie dla admina)
func (h *Handler) GetForFacilityOwner(c *gin.Context) {
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

	role, _ := c.Get("role")
	if role != nil && role.(string) == "ADMIN" {
		reservations, err := h.service.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": reservations})
		return
	}

	reservations, err := h.service.GetByFacilityOwnerID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reservations})
}

// UpdateStatus – pozwala właścicielowi obiektu potwierdzić lub odrzucić oczekującą rezerwację
func (h *Handler) UpdateStatus(c *gin.Context) {
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

	reservationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.ReservationStatus(body.Status)
	if status != models.StatusConfirmed && status != models.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be CONFIRMED or CANCELLED"})
		return
	}

	if err := h.service.UpdateReservationStatus(reservationID, userID, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation status updated"})
}
