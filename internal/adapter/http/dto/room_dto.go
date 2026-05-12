package dto

// CreateRoomRequest is the HTTP request body for POST /api/v1/rooms.
type CreateRoomRequest struct {
	OrganizerName string `json:"organizerName"`
}

// CreateRoomResponse is the HTTP response for POST /api/v1/rooms.
type CreateRoomResponse struct {
	RoomID         string `json:"roomId"`
	InviteURL      string `json:"inviteUrl"`
	OrganizerToken string `json:"token"`
}

// JoinRoomRequest is the HTTP request body for POST /api/v1/rooms/:id/join.
type JoinRoomRequest struct {
	ParticipantName string `json:"participantName"`
}

// JoinRoomResponse is the HTTP response for POST /api/v1/rooms/:id/join.
type JoinRoomResponse struct {
	ParticipantID       string              `json:"participantId"`
	ParticipantToken    string              `json:"token"`
	CurrentParticipants []ParticipantInfo   `json:"participants"`
}

// ParticipantInfo is a lightweight participant representation for HTTP responses.
type ParticipantInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// SendMessageRequest is the HTTP request body for POST /api/v1/rooms/:id/chat.
type SendMessageRequest struct {
	Text string `json:"text"`
}

// ErrorResponse is the standard error envelope for all HTTP errors.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
