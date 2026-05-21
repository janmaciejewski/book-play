package user

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

type UpdateRoleDTO struct {
	Role string `json:"role" binding:"required"`
}

func (h *Handler) GetAll(c *gin.Context) {
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
	type R struct {
		ID        string  `json:"id"`
		Email     string  `json:"email"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Phone     *string `json:"phone,omitempty"`
		Role      string  `json:"role"`
		IsActive  bool    `json:"is_active"`
	}
	resp := make([]R, len(users))
	for i, u := range users {
		resp[i] = R{ID: u.ID.String(), Email: u.Email, FirstName: u.FirstName, LastName: u.LastName, Phone: u.Phone, Role: string(u.Role), IsActive: u.IsActive}
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) GetProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))
	paramID, err := uuid.Parse(c.Param("id"))
	if err != nil || userID != paramID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}
	var dto UpdateProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.UpdateProfile(userID, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	userRole, _ := c.Get("role")
	if userRole != models.RoleAdmin && userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}
	id, _ := uuid.Parse(c.Param("id"))
	var dto UpdateRoleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	valid := map[string]bool{"ADMIN": true, "FACILITY_OWNER": true, "PLAYER": true}
	if !valid[dto.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	h.service.UpdateRole(id, dto.Role)
	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))
	paramID, err := uuid.Parse(c.Param("id"))
	if err != nil || userID != paramID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own avatar"})
		return
	}
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No avatar file provided"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read avatar file"})
		return
	}
	if len(fileData) > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file too large (max 5MB)"})
		return
	}
	avatarPath, err := h.service.SaveAvatar(userID, fileData, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar": avatarPath})
}

func (h *Handler) Delete(c *gin.Context) {
	userRole, _ := c.Get("role")
	if userRole != models.RoleAdmin && userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}
	id, _ := uuid.Parse(c.Param("id"))
	h.service.Delete(id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
