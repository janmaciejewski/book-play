package user

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
// @Summary Get all users
// @Description Get list of all users (admin only)
// @Tags users
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /users [get]
func (h *Handler) GetAll(c *gin.Context) {
	// Check if user is admin
	userRole, _ := c.Get("role")
	if userRole != models.RoleAdmin && userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Remove password hashes from response
	type UserResponse struct {
		ID        string  `json:"id"`
		Email     string  `json:"email"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Phone     *string `json:"phone,omitempty"`
		Role      string  `json:"role"`
		IsActive  bool    `json:"is_active"`
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateRoleDTO for updating user role
type UpdateRoleDTO struct {
	Role string `json:"role" binding:"required"`
}

// UpdateRole godoc
// @Summary Update user role
// @Description Update a user's role (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateRoleDTO true "New role"
// @Success 200 {object} map[string]string
// @Router /users/{id}/role [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	// Check if user is admin
	userRole, _ := c.Get("role")
	if userRole != models.RoleAdmin && userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var dto UpdateRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"ADMIN":          true,
		"FACILITY_OWNER": true,
		"PLAYER":         true,
	}
	if !validRoles[dto.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	if err := h.service.UpdateRole(id, dto.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	// Check if user is admin
	userRole, _ := c.Get("role")
	if userRole != models.RoleAdmin && userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
