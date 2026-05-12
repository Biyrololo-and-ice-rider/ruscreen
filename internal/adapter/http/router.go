package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourorg/meet/internal/adapter/http/handler"
)

// NewRouter wires all HTTP routes.
func NewRouter(rooms *handler.RoomHandler, chat *handler.ChatHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/rooms", rooms.Create)
		r.Post("/rooms/{id}/join", rooms.Join)
		r.Delete("/rooms/{id}", rooms.Close)
		r.Delete("/rooms/{id}/leave", rooms.Leave)
		r.Get("/rooms/{id}/chat", chat.GetHistory)
		r.Post("/rooms/{id}/chat", chat.Send)
	})

	return r
}
