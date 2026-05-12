package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/yourorg/meet/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// TODO: wire up dependencies
	// tokenGen  := token.NewJWTGenerator(cfg.JWTSecret)
	// redisClient := redis.NewClient(cfg.RedisURL)
	// roomRepo  := redisinfra.NewRoomRepository(redisClient)
	// chatRepo  := redisinfra.NewChatRepository(redisClient)
	// hub       := wsinfra.NewHub()
	// go hub.Run()
	//
	// createRoom  := usecase.NewCreateRoom(roomRepo, tokenGen)
	// joinRoom    := usecase.NewJoinRoom(roomRepo, tokenGen)
	// ...
	//
	// r := router.New(createRoom, joinRoom, ...)
	// wsHandler := wsadapter.NewHandler(hub, ...)

	slog.Info("starting meet server", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
