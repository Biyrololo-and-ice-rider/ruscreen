package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/meet/internal/adapter/http/dto"
	"github.com/yourorg/meet/internal/usecase"
)

// ChatHandler handles HTTP requests for chat operations.
type ChatHandler struct {
	send       *usecase.SendChatMessage
	getHistory *usecase.GetChatHistory
}

func NewChatHandler(send *usecase.SendChatMessage, get *usecase.GetChatHistory) *ChatHandler {
	return &ChatHandler{send: send, getHistory: get}
}

// POST /api/v1/rooms/:id/chat
func (h *ChatHandler) Send(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	result, err := h.send.Execute(r.Context(), usecase.SendChatMessageRequest{
		RoomID:           chi.URLParam(r, "id"),
		ParticipantToken: token,
		Text:             req.Text,
	})
	if err != nil {
		writeFromDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"messageId": result.MessageID,
		"sentAt":    result.SentAt,
	})
}

// GET /api/v1/rooms/:id/chat
func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	result, err := h.getHistory.Execute(r.Context(), usecase.GetChatHistoryRequest{
		RoomID:           chi.URLParam(r, "id"),
		ParticipantToken: token,
	})
	if err != nil {
		writeFromDomainErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"messages": result.Messages})
}
