package port

import (
	"context"

	"github.com/yourorg/meet/internal/domain"
)

// ChatRepository is the outgoing port for chat message persistence.
// Implementations live in internal/infra/redis/.
type ChatRepository interface {
	Save(ctx context.Context, msg domain.ChatMessage) error
	FindByRoom(ctx context.Context, roomID string) ([]domain.ChatMessage, error)
}
