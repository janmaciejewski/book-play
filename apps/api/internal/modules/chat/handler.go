package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	noRedis bool
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func NewNoRedisHandler() *Handler {
	return &Handler{noRedis: true}
}

type SendMessageDTO struct {
	Text     string `json:"text" binding:"required"`
	UserName string `json:"user_name"`
}

func (h *Handler) GetMessages(c *gin.Context) {
	if h.noRedis {
		c.JSON(http.StatusOK, gin.H{"data": []*Message{}})
		return
	}
	teamID := c.Param("id")
	since := c.Query("since")

	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")

	_ = userID
	_ = userName

	messages, err := h.service.GetMessages(teamID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	if messages == nil {
		messages = []*Message{}
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

func (h *Handler) SendMessage(c *gin.Context) {
	teamID := c.Param("id")

	userIDStr, _ := c.Get("userID")
	userID := userIDStr.(string)

	var dto SendMessageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if dto.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message text is required"})
		return
	}

	userName := dto.UserName
	if userName == "" {
		userName = userID
	}

	role, _ := c.Get("role")
	userRole := ""
	if role != nil {
		userRole = role.(string)
	}

	msg, err := h.service.SendMessage(teamID, userID, userName, userRole, dto.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": msg})
}
