package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/yourorg/meet/internal/domain"
)

// Client represents a connected WebSocket participant.
type Client struct {
	ParticipantID string
	RoomID        string
	conn          *websocket.Conn
	send          chan []byte
}

// Hub implements port.SignalingGateway via an in-memory goroutine fan-out map.
type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[string]*Client // roomID → participantID → Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[string]*Client),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
	}
}

// Run starts the hub's event loop. Call in a goroutine: go hub.Run().
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.rooms[c.RoomID] == nil {
				h.rooms[c.RoomID] = make(map[string]*Client)
			}
			h.rooms[c.RoomID][c.ParticipantID] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[c.RoomID]; ok {
				delete(room, c.ParticipantID)
				if len(room) == 0 {
					delete(h.rooms, c.RoomID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Register(c *Client)   { h.register <- c }
func (h *Hub) Unregister(c *Client) { h.unregister <- c }

// broadcast sends a JSON payload to all clients in a room.
func (h *Hub) broadcast(roomID string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.rooms[roomID] {
		select {
		case c.send <- data:
		default: // slow client — skip
		}
	}
	return nil
}

// --- port.SignalingGateway implementation ---

func (h *Hub) BroadcastRoomClosed(_ context.Context, roomID string) error {
	return h.broadcast(roomID, map[string]string{"type": "room_closed", "roomId": roomID})
}

func (h *Hub) BroadcastParticipantLeft(_ context.Context, roomID, participantID string) error {
	return h.broadcast(roomID, map[string]string{"type": "participant_left", "participantId": participantID})
}

func (h *Hub) BroadcastMediaState(_ context.Context, roomID, participantID string, camera, mic bool) error {
	return h.broadcast(roomID, map[string]any{
		"type":          "media_state",
		"participantId": participantID,
		"camera":        camera,
		"mic":           mic,
	})
}

func (h *Hub) BroadcastChatMessage(_ context.Context, roomID string, msg domain.ChatMessage) error {
	return h.broadcast(roomID, map[string]any{
		"type":        "chat_message",
		"id":          msg.ID,
		"senderName":  msg.SenderName,
		"text":        msg.Text,
		"sentAt":      msg.SentAt,
	})
}

func (h *Hub) BroadcastScreenShareState(_ context.Context, roomID, participantID string, active bool) error {
	return h.broadcast(roomID, map[string]any{
		"type":          "screen_share_state",
		"participantId": participantID,
		"active":        active,
	})
}

func (h *Hub) ForwardSDP(_ context.Context, roomID, _, toID, sdp string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room, ok := h.rooms[roomID]
	if !ok {
		return nil
	}
	target, ok := room[toID]
	if !ok {
		return nil
	}
	data, _ := json.Marshal(map[string]string{"type": "sdp", "sdp": sdp})
	select {
	case target.send <- data:
	default:
	}
	return nil
}
