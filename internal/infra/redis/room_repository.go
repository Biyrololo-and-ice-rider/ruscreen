package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/meet/internal/domain"
)

const roomTTL = 24 * time.Hour

// RoomRepository implements port.RoomRepository using Redis.
type RoomRepository struct {
	client *redis.Client
}

func NewRoomRepository(client *redis.Client) *RoomRepository {
	return &RoomRepository{client: client}
}

func (r *RoomRepository) Save(ctx context.Context, room domain.Room) error {
	data, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("marshal room: %w", err)
	}
	key := "room:" + room.ID
	return r.client.Set(ctx, key, data, roomTTL).Err()
}

func (r *RoomRepository) FindByID(ctx context.Context, roomID string) (domain.Room, error) {
	key := "room:" + roomID
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return domain.Room{}, domain.ErrRoomNotFound
	}
	if err != nil {
		return domain.Room{}, fmt.Errorf("redis get: %w", err)
	}
	var room domain.Room
	if err := json.Unmarshal(data, &room); err != nil {
		return domain.Room{}, fmt.Errorf("unmarshal room: %w", err)
	}
	return room, nil
}

func (r *RoomRepository) UpdateStatus(ctx context.Context, roomID string, status domain.RoomStatus) error {
	room, err := r.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	room.Status = status
	return r.Save(ctx, room)
}

func (r *RoomRepository) AddParticipant(ctx context.Context, roomID string, p domain.Participant) error {
	room, err := r.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	room.Participants = append(room.Participants, p)
	return r.Save(ctx, room)
}

func (r *RoomRepository) RemoveParticipant(ctx context.Context, roomID string, participantID string) error {
	room, err := r.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	filtered := room.Participants[:0]
	for _, p := range room.Participants {
		if p.ID != participantID {
			filtered = append(filtered, p)
		}
	}
	room.Participants = filtered
	return r.Save(ctx, room)
}

func (r *RoomRepository) SetScreenSharer(ctx context.Context, roomID string, participantID *string) error {
	room, err := r.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	room.ScreenSharerID = participantID
	return r.Save(ctx, room)
}
