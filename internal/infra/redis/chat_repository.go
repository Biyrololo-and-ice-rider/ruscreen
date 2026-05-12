package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/meet/internal/domain"
)

const chatTTL = 24 * time.Hour

// ChatRepository implements port.ChatRepository using a Redis List.
type ChatRepository struct {
	client *redis.Client
}

func NewChatRepository(client *redis.Client) *ChatRepository {
	return &ChatRepository{client: client}
}

func (r *ChatRepository) Save(ctx context.Context, msg domain.ChatMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	key := "chat:" + msg.RoomID + ":messages"
	pipe := r.client.Pipeline()
	pipe.RPush(ctx, key, data)
	pipe.Expire(ctx, key, chatTTL)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *ChatRepository) FindByRoom(ctx context.Context, roomID string) ([]domain.ChatMessage, error) {
	key := "chat:" + roomID + ":messages"
	items, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange: %w", err)
	}
	msgs := make([]domain.ChatMessage, 0, len(items))
	for _, item := range items {
		var msg domain.ChatMessage
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			continue // skip corrupt entries
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
