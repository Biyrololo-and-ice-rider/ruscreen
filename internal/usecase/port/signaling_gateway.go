package port

import (
	"context"

	"github.com/yourorg/meet/internal/domain"
)

// SignalingGateway is the outgoing port for real-time WebSocket events.
// The implementation is the in-memory WebSocket Hub (internal/infra/ws/hub.go).
type SignalingGateway interface {
	// BroadcastRoomClosed notifies all participants that the room has ended.
	BroadcastRoomClosed(ctx context.Context, roomID string) error

	// BroadcastParticipantLeft notifies remaining participants of a departure.
	BroadcastParticipantLeft(ctx context.Context, roomID, participantID string) error

	// BroadcastMediaState notifies participants of camera/mic state change.
	BroadcastMediaState(ctx context.Context, roomID, participantID string, camera, mic bool) error

	// BroadcastChatMessage fans out a chat message to all room participants.
	BroadcastChatMessage(ctx context.Context, roomID string, msg domain.ChatMessage) error

	// BroadcastScreenShareState notifies participants of screen share start/stop.
	BroadcastScreenShareState(ctx context.Context, roomID, participantID string, active bool) error

	// ForwardSDP relays a WebRTC SDP offer/answer from one peer to another.
	ForwardSDP(ctx context.Context, roomID, fromID, toID, sdp string) error
}
