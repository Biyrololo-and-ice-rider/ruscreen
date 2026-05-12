package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase/port"
)

// ----- CloseRoom -----

type CloseRoomRequest struct {
	RoomID         string
	OrganizerToken string
}

type CloseRoom struct {
	rooms     port.RoomRepository
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewCloseRoom(rooms port.RoomRepository, tokens port.TokenGenerator, sig port.SignalingGateway) *CloseRoom {
	return &CloseRoom{rooms: rooms, tokens: tokens, signaling: sig}
}

func (uc *CloseRoom) Execute(ctx context.Context, req CloseRoomRequest) error {
	_, _, role, err := uc.tokens.Verify(req.OrganizerToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	if role != domain.RoleOrganizer {
		return domain.ErrForbidden
	}

	room, err := uc.rooms.FindByID(ctx, req.RoomID)
	if err != nil {
		return err
	}
	if err := room.Close(time.Now().UTC()); err != nil {
		return err
	}
	if err := uc.rooms.UpdateStatus(ctx, req.RoomID, domain.RoomStatusClosed); err != nil {
		return fmt.Errorf("update room status: %w", err)
	}
	return uc.signaling.BroadcastRoomClosed(ctx, req.RoomID)
}

// ----- LeaveRoom -----

type LeaveRoomRequest struct {
	RoomID           string
	ParticipantToken string
}

type LeaveRoom struct {
	rooms     port.RoomRepository
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewLeaveRoom(rooms port.RoomRepository, tokens port.TokenGenerator, sig port.SignalingGateway) *LeaveRoom {
	return &LeaveRoom{rooms: rooms, tokens: tokens, signaling: sig}
}

func (uc *LeaveRoom) Execute(ctx context.Context, req LeaveRoomRequest) error {
	_, participantID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	if err := uc.rooms.RemoveParticipant(ctx, req.RoomID, participantID); err != nil {
		return fmt.Errorf("remove participant: %w", err)
	}
	return uc.signaling.BroadcastParticipantLeft(ctx, req.RoomID, participantID)
}

// ----- ToggleCamera -----

type ToggleMediaRequest struct {
	RoomID           string
	ParticipantToken string
	Enabled          bool
}

type ToggleCamera struct {
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewToggleCamera(tokens port.TokenGenerator, sig port.SignalingGateway) *ToggleCamera {
	return &ToggleCamera{tokens: tokens, signaling: sig}
}

func (uc *ToggleCamera) Execute(ctx context.Context, req ToggleMediaRequest) error {
	_, participantID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	return uc.signaling.BroadcastMediaState(ctx, req.RoomID, participantID, req.Enabled, true)
}

// ----- ToggleMic -----

type ToggleMic struct {
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewToggleMic(tokens port.TokenGenerator, sig port.SignalingGateway) *ToggleMic {
	return &ToggleMic{tokens: tokens, signaling: sig}
}

func (uc *ToggleMic) Execute(ctx context.Context, req ToggleMediaRequest) error {
	_, participantID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return domain.ErrInvalidToken
	}
	return uc.signaling.BroadcastMediaState(ctx, req.RoomID, participantID, true, req.Enabled)
}
