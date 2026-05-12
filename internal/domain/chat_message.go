package domain

import "time"

// ChatMessage represents a single text message sent in a room's chat.
type ChatMessage struct {
	ID         string
	RoomID     string
	SenderID   string
	SenderName string // denormalized for history display
	Text       string
	SentAt     time.Time
}

// NewChatMessage creates a ChatMessage, validating text invariants.
func NewChatMessage(id, roomID, senderID, senderName, text string, sentAt time.Time) (ChatMessage, error) {
	if text == "" {
		return ChatMessage{}, ErrEmptyMessage
	}
	if len([]rune(text)) > 2000 {
		return ChatMessage{}, ErrMessageTooLong
	}
	return ChatMessage{
		ID:         id,
		RoomID:     roomID,
		SenderID:   senderID,
		SenderName: senderName,
		Text:       text,
		SentAt:     sentAt,
	}, nil
}
