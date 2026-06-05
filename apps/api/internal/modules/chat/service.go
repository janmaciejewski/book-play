package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	ID        string `json:"id"`
	TeamID    string `json:"teamId"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserRole  string `json:"userRole,omitempty"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type Service struct {
	redis *redis.Client
}

func NewService(redis *redis.Client) *Service {
	return &Service{redis: redis}
}

func chatKey(teamID string) string {
	return "team:chat:" + teamID
}

func (s *Service) SendMessage(teamID, userID, userName, userRole, text string) (*Message, error) {
	ctx := context.Background()

	msg := &Message{
		ID:        uuid.New().String(),
		TeamID:    teamID,
		UserID:    userID,
		UserName:  userName,
		UserRole:  userRole,
		Text:      text,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	score := float64(time.Now().UnixMilli())

	if err := s.redis.ZAdd(ctx, chatKey(teamID), redis.Z{
		Score:  score,
		Member: string(data),
	}).Err(); err != nil {
		return nil, fmt.Errorf("failed to store message: %w", err)
	}

	// Ustawia 24h TTL na secie – odświeża się przy każdej wiadomości
	s.redis.Expire(ctx, chatKey(teamID), 24*time.Hour)

	return msg, nil
}

func (s *Service) GetMessages(teamID string, since string) ([]*Message, error) {
	ctx := context.Background()

	minScore := "-inf"
	if since != "" {
		t, err := time.Parse(time.RFC3339, since)
		if err == nil {
			minScore = fmt.Sprintf("%d", t.UnixMilli()+1)
		}
	}

	results, err := s.redis.ZRangeByScore(ctx, chatKey(teamID), &redis.ZRangeBy{
		Min: minScore,
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	messages := make([]*Message, 0, len(results))
	for _, raw := range results {
		var msg Message
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
