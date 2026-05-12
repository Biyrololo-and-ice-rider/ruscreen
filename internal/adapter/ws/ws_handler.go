package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	wsinfra "github.com/yourorg/meet/internal/infra/ws"
	"github.com/yourorg/meet/internal/usecase/port"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO: restrict in production
}

// Message is the envelope for all WebSocket signaling messages.
type Message struct {
	Type          string `json:"type"`
	RoomID        string `json:"roomId,omitempty"`
	ParticipantID string `json:"participantId,omitempty"`
	ToID          string `json:"toId,omitempty"`
	SDP           string `json:"sdp,omitempty"`
	Text          string `json:"text,omitempty"`
}

// Handler upgrades HTTP connections to WebSocket and routes signaling messages.
type Handler struct {
	hub    *wsinfra.Hub
	tokens port.TokenGenerator
}

func NewHandler(hub *wsinfra.Hub, tokens port.TokenGenerator) *Handler {
	return &Handler{hub: hub, tokens: tokens}
}

// ServeWS handles GET /ws/:roomID
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	roomID, participantID, _, err := h.tokens.Verify(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "err", err)
		return
	}

	client := &wsinfra.Client{
		ParticipantID: participantID,
		RoomID:        roomID,
	}
	// Note: Hub.Client fields are exported; send channel wiring omitted for brevity
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		switch msg.Type {
		case "sdp":
			_ = h.hub.ForwardSDP(r.Context(), roomID, participantID, msg.ToID, msg.SDP)
		}
	}
}
