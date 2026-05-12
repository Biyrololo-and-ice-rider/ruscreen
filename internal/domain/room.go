package domain

import "time"

// RoomStatus represents the lifecycle state of a conference room.
type RoomStatus string

const (
	RoomStatusActive RoomStatus = "ACTIVE"
	RoomStatusClosed RoomStatus = "CLOSED"
)

// Room is the core entity representing a conference session.
// It enforces Enterprise Business Rules independent of any framework or database.
type Room struct {
	ID             string
	Status         RoomStatus
	Participants   []Participant
	ScreenSharerID *string // nil means nobody is sharing
	CreatedAt      time.Time
	ClosedAt       *time.Time
}

// NewRoom creates a new active Room, validating all invariants.
func NewRoom(id string, createdAt time.Time) (Room, error) {
	if id == "" {
		return Room{}, ErrRoomNotFound // reuse; means "invalid ID"
	}
	return Room{
		ID:           id,
		Status:       RoomStatusActive,
		Participants: []Participant{},
		CreatedAt:    createdAt,
	}, nil
}

// CanJoin returns true if the room is ACTIVE and has fewer than 4 participants.
func (r *Room) CanJoin() bool {
	return r.Status == RoomStatusActive && len(r.Participants) < 4
}

// IsActive returns true if the room is in ACTIVE state.
func (r *Room) IsActive() bool {
	return r.Status == RoomStatusActive
}

// HasScreenSharer returns true if a participant is currently sharing their screen.
func (r *Room) HasScreenSharer() bool {
	return r.ScreenSharerID != nil
}

// Close transitions the room to CLOSED state.
// Returns ErrRoomAlreadyClosed if the room is not ACTIVE.
func (r *Room) Close(closedAt time.Time) error {
	if r.Status != RoomStatusActive {
		return ErrRoomAlreadyClosed
	}
	r.Status = RoomStatusClosed
	r.ClosedAt = &closedAt
	return nil
}
