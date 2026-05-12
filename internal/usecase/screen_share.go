package usecase

import (
	"context"

	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase/port"
)

type ScreenShareRequest struct {
	RoomID           string
	ParticipantToken string
}

// StartScreenShare enables screen sharing for a participant.
// Only one participant may share at a time.
type StartScreenShare struct {
	rooms     port.RoomRepository
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewStartScreenShare(rooms port.RoomRepository, tokens port.TokenGenerator, sig port.SignalingGateway) *StartScreenShare {
	return &StartScreenShare{rooms: rooms, tokens: tokens, signaling: sig}
}

func (uc *StartScreenShare) Execute(ctx context.Context, req ScreenShareRequest) error {
	_, participantID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return domain.ErrInvalidToken
	}

	room, err := uc.rooms.FindByID(ctx, req.RoomID)
	if err != nil {
		return err
	}
	if room.HasScreenSharer() {
		return domain.ErrScreenShareActive
	}

	if err := uc.rooms.SetScreenSharer(ctx, req.RoomID, &participantID); err != nil {
		return err
	}
	return uc.signaling.BroadcastScreenShareState(ctx, req.RoomID, participantID, true)
}

// StopScreenShare ends the current screen share session.
type StopScreenShare struct {
	rooms     port.RoomRepository
	tokens    port.TokenGenerator
	signaling port.SignalingGateway
}

func NewStopScreenShare(rooms port.RoomRepository, tokens port.TokenGenerator, sig port.SignalingGateway) *StopScreenShare {
	return &StopScreenShare{rooms: rooms, tokens: tokens, signaling: sig}
}

func (uc *StopScreenShare) Execute(ctx context.Context, req ScreenShareRequest) error {
	_, participantID, _, err := uc.tokens.Verify(req.ParticipantToken)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if err := uc.rooms.SetScreenSharer(ctx, req.RoomID, nil); err != nil {
		return err
	}
	return uc.signaling.BroadcastScreenShareState(ctx, req.RoomID, participantID, false)
}
