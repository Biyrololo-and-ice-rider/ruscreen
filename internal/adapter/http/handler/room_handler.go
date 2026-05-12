package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/meet/internal/adapter/http/dto"
	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase"
)

// RoomHandler handles HTTP requests for room lifecycle operations.
type RoomHandler struct {
	createRoom *usecase.CreateRoom
	joinRoom   *usecase.JoinRoom
	closeRoom  *usecase.CloseRoom
	leaveRoom  *usecase.LeaveRoom
}

func NewRoomHandler(
	create *usecase.CreateRoom,
	join *usecase.JoinRoom,
	close *usecase.CloseRoom,
	leave *usecase.LeaveRoom,
) *RoomHandler {
	return &RoomHandler{createRoom: create, joinRoom: join, closeRoom: close, leaveRoom: leave}
}

// POST /api/v1/rooms
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	result, err := h.createRoom.Execute(r.Context(), usecase.CreateRoomRequest{
		OrganizerName: req.OrganizerName,
	})
	if err != nil {
		writeFromDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateRoomResponse{
		RoomID:         result.RoomID,
		InviteURL:      result.InviteURL,
		OrganizerToken: result.OrganizerToken,
	})
}

// POST /api/v1/rooms/:id/join
func (h *RoomHandler) Join(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	var req dto.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	result, err := h.joinRoom.Execute(r.Context(), usecase.JoinRoomRequest{
		RoomID:          roomID,
		ParticipantName: req.ParticipantName,
	})
	if err != nil {
		writeFromDomainErr(w, err)
		return
	}

	participants := make([]dto.ParticipantInfo, 0, len(result.CurrentParticipants))
	for _, p := range result.CurrentParticipants {
		participants = append(participants, dto.ParticipantInfo{ID: p.ID, Name: p.Name, Role: p.Role})
	}

	writeJSON(w, http.StatusOK, dto.JoinRoomResponse{
		ParticipantID:       result.ParticipantID,
		ParticipantToken:    result.ParticipantToken,
		CurrentParticipants: participants,
	})
}

// DELETE /api/v1/rooms/:id
func (h *RoomHandler) Close(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	token := r.Header.Get("Authorization")

	if err := h.closeRoom.Execute(r.Context(), usecase.CloseRoomRequest{
		RoomID:         roomID,
		OrganizerToken: token,
	}); err != nil {
		writeFromDomainErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/rooms/:id/leave
func (h *RoomHandler) Leave(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	token := r.Header.Get("Authorization")

	if err := h.leaveRoom.Execute(r.Context(), usecase.LeaveRoomRequest{
		RoomID:           roomID,
		ParticipantToken: token,
	}); err != nil {
		writeFromDomainErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, dto.ErrorResponse{Code: code, Message: msg})
}

func writeFromDomainErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrRoomNotFound):
		writeError(w, http.StatusNotFound, "ROOM_NOT_FOUND", err.Error())
	case errors.Is(err, domain.ErrRoomClosed):
		writeError(w, http.StatusGone, "ROOM_CLOSED", err.Error())
	case errors.Is(err, domain.ErrRoomFull):
		writeError(w, http.StatusConflict, "ROOM_FULL", err.Error())
	case errors.Is(err, domain.ErrEmptyName), errors.Is(err, domain.ErrNameTooLong):
		writeError(w, http.StatusBadRequest, "INVALID_NAME", err.Error())
	case errors.Is(err, domain.ErrForbidden):
		writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, domain.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}
