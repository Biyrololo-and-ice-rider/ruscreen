package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase/port"
)

// CreateRoomRequest is the input model for the CreateRoom use case.
type CreateRoomRequest struct {
	OrganizerName string
}

// CreateRoomResponse is the output model for the CreateRoom use case.
type CreateRoomResponse struct {
	RoomID        string
	InviteURL     string
	OrganizerToken string
}

// CreateRoom orchestrates creating a new conference room.
type CreateRoom struct {
	rooms  port.RoomRepository
	tokens port.TokenGenerator
}

func NewCreateRoom(rooms port.RoomRepository, tokens port.TokenGenerator) *CreateRoom {
	return &CreateRoom{rooms: rooms, tokens: tokens}
}

func (uc *CreateRoom) Execute(ctx context.Context, req CreateRoomRequest) (CreateRoomResponse, error) {
	if req.OrganizerName == "" {
		return CreateRoomResponse{}, domain.ErrEmptyName
	}

	roomID := uuid.NewString()[:8] // short ID for URLs
	now := time.Now().UTC()

	room, err := domain.NewRoom(roomID, now)
	if err != nil {
		return CreateRoomResponse{}, err
	}

	organizerID := uuid.NewString()
	organizer, err := domain.NewParticipant(organizerID, req.OrganizerName, domain.RoleOrganizer, now)
	if err != nil {
		return CreateRoomResponse{}, err
	}
	room.Participants = append(room.Participants, organizer)

	if err := uc.rooms.Save(ctx, room); err != nil {
		return CreateRoomResponse{}, fmt.Errorf("save room: %w", err)
	}

	token, err := uc.tokens.NewToken(roomID, organizerID, domain.RoleOrganizer)
	if err != nil {
		return CreateRoomResponse{}, fmt.Errorf("generate token: %w", err)
	}

	return CreateRoomResponse{
		RoomID:        roomID,
		InviteURL:     "/room/" + roomID,
		OrganizerToken: token,
	}, nil
}
