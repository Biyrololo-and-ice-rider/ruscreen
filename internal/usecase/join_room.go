package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase/port"
)

// JoinRoomRequest is the input model for the JoinRoom use case.
type JoinRoomRequest struct {
	RoomID          string
	ParticipantName string
}

// ParticipantInfo is a lightweight view of a participant for the response.
type ParticipantInfo struct {
	ID   string
	Name string
	Role string
}

// JoinRoomResponse is the output model for the JoinRoom use case.
type JoinRoomResponse struct {
	ParticipantID    string
	ParticipantToken string
	CurrentParticipants []ParticipantInfo
}

// JoinRoom orchestrates adding a participant to an existing room.
type JoinRoom struct {
	rooms  port.RoomRepository
	tokens port.TokenGenerator
}

func NewJoinRoom(rooms port.RoomRepository, tokens port.TokenGenerator) *JoinRoom {
	return &JoinRoom{rooms: rooms, tokens: tokens}
}

func (uc *JoinRoom) Execute(ctx context.Context, req JoinRoomRequest) (JoinRoomResponse, error) {
	if req.ParticipantName == "" {
		return JoinRoomResponse{}, domain.ErrEmptyName
	}

	room, err := uc.rooms.FindByID(ctx, req.RoomID)
	if err != nil {
		return JoinRoomResponse{}, err
	}

	if !room.IsActive() {
		return JoinRoomResponse{}, domain.ErrRoomClosed
	}
	if !room.CanJoin() {
		return JoinRoomResponse{}, domain.ErrRoomFull
	}

	now := time.Now().UTC()
	participantID := uuid.NewString()
	participant, err := domain.NewParticipant(participantID, req.ParticipantName, domain.RoleGuest, now)
	if err != nil {
		return JoinRoomResponse{}, err
	}

	if err := uc.rooms.AddParticipant(ctx, req.RoomID, participant); err != nil {
		return JoinRoomResponse{}, fmt.Errorf("add participant: %w", err)
	}

	token, err := uc.tokens.NewToken(req.RoomID, participantID, domain.RoleGuest)
	if err != nil {
		return JoinRoomResponse{}, fmt.Errorf("generate token: %w", err)
	}

	current := make([]ParticipantInfo, 0, len(room.Participants)+1)
	for _, p := range room.Participants {
		current = append(current, ParticipantInfo{ID: p.ID, Name: p.Name, Role: string(p.Role)})
	}
	current = append(current, ParticipantInfo{ID: participantID, Name: req.ParticipantName, Role: string(domain.RoleGuest)})

	return JoinRoomResponse{
		ParticipantID:       participantID,
		ParticipantToken:    token,
		CurrentParticipants: current,
	}, nil
}
