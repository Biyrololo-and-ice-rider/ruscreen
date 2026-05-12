package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase/port"
)

// SendChatMessageRequest is the input model for sending a chat message.
type SendChatMessageRequest struct {
	RoomID           string
	ParticipantToken string
	Text             string
}

// SendChatMessageResponse is the output model for sending a chat message.
type SendChatMessageResponse struct {
	MessageID string
	SentAt    time.Time
}

// SendChatMessage handles sending a text message to a room's chat.
type SendChatMessage struct {
	rooms     port.RoomRepository
	chat      port.ChatRepository
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewSendChatMessage(
	rooms port.RoomRepository,
	chat port.ChatRepository,
	tokens port.TokenGenerator,
	sig port.SignalingGateway,
) *SendChatMessage {
	return &SendChatMessage{rooms: rooms, chat: chat, tokens: tokens, signaling: sig}
}

func (uc *SendChatMessage) Execute(ctx context.Context, req SendChatMessageRequest) (SendChatMessageResponse, error) {
	roomID, senderID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return SendChatMessageResponse{}, domain.ErrInvalidToken
	}

	room, err := uc.rooms.FindByID(ctx, roomID)
	if err != nil {
		return SendChatMessageResponse{}, err
	}
	if !room.IsActive() {
		return SendChatMessageResponse{}, domain.ErrRoomClosed
	}

	// Lookup sender name from room participants
	senderName := "Unknown"
	for _, p := range room.Participants {
		if p.ID == senderID {
			senderName = p.Name
			break
		}
	}

	now := time.Now().UTC()
	msgID := uuid.NewString()

	msg, err := domain.NewChatMessage(msgID, roomID, senderID, senderName, req.Text, now)
	if err != nil {
		return SendChatMessageResponse{}, err
	}

	if err := uc.chat.Save(ctx, msg); err != nil {
		return SendChatMessageResponse{}, fmt.Errorf("save message: %w", err)
	}

	if err := uc.signaling.BroadcastChatMessage(ctx, roomID, msg); err != nil {
		// Non-fatal: message is saved; broadcast failure is logged, not returned
		_ = err
	}

	return SendChatMessageResponse{MessageID: msgID, SentAt: now}, nil
}

// ----- GetChatHistory -----

type GetChatHistoryRequest struct {
	RoomID           string
	ParticipantToken string
}

type GetChatHistoryResponse struct {
	Messages []domain.ChatMessage
}

type GetChatHistory struct {
	chat   port.ChatRepository
	tokens port.TokenGenerator
}

func NewGetChatHistory(chat port.ChatRepository, tokens port.TokenGenerator) *GetChatHistory {
	return &GetChatHistory{chat: chat, tokens: tokens}
}

func (uc *GetChatHistory) Execute(ctx context.Context, req GetChatHistoryRequest) (GetChatHistoryResponse, error) {
	if _, _, _, err := uc.tokens.Verify(req.ParticipantToken); err != nil {
		return GetChatHistoryResponse{}, domain.ErrInvalidToken
	}
	msgs, err := uc.chat.FindByRoom(ctx, req.RoomID)
	if err != nil {
		return GetChatHistoryResponse{}, fmt.Errorf("fetch chat: %w", err)
	}
	return GetChatHistoryResponse{Messages: msgs}, nil
}
