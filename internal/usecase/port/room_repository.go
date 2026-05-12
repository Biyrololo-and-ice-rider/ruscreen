package port

import (
	"context"

	"github.com/yourorg/meet/internal/domain"
)

// RoomRepository is the outgoing port for persisting and querying rooms.
// Implementations live in internal/infra/redis/.
type RoomRepository interface {
	Save(ctx context.Context, room domain.Room) error
	FindByID(ctx context.Context, roomID string) (domain.Room, error)
	UpdateStatus(ctx context.Context, roomID string, status domain.RoomStatus) error
	AddParticipant(ctx context.Context, roomID string, p domain.Participant) error
	RemoveParticipant(ctx context.Context, roomID string, participantID string) error
	SetScreenSharer(ctx context.Context, roomID string, participantID *string) error
}
