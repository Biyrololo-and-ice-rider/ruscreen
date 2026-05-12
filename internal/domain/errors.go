package domain

import "errors"

var (
	// Room errors
	ErrRoomNotFound      = errors.New("room not found")
	ErrRoomClosed        = errors.New("room is closed")
	ErrRoomFull          = errors.New("room is full (max 4 participants)")
	ErrRoomAlreadyClosed = errors.New("room is already closed")

	// Participant errors
	ErrEmptyName  = errors.New("name must not be empty")
	ErrNameTooLong = errors.New("name must not exceed 50 characters")

	// Chat errors
	ErrEmptyMessage   = errors.New("message must not be empty")
	ErrMessageTooLong = errors.New("message must not exceed 2000 characters")

	// Auth errors
	ErrForbidden    = errors.New("forbidden: insufficient role")
	ErrInvalidToken = errors.New("invalid or expired token")

	// Screen share errors
	ErrScreenShareActive = errors.New("screen share already active by another participant")
)
