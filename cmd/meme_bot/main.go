package main

import (
	"log/slog"
	"meme_bot/internal/bot"
	"meme_bot/internal/handler"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	token := os.Getenv("TOKEN")
	chatID := os.Getenv("CHAT_ID")
	slog.Info("starting app", "port", port, "token", token, "chatID", chatID)

	if token == "" {
		slog.Error("token is empty")
		os.Exit(1)
	}
	if chatID == "" {
		slog.Error("chatID is empty")
		os.Exit(1)
	}
	if port == "" {
		slog.Error("PORT is empty")
		os.Exit(1)
	}

	b := bot.NewBot(token)
	h := handler.NewHandler(chatID, b)

	router := http.NewServeMux()
	router.HandleFunc("/webhook", h.HandleMessage)

	srv := NewServer(port, router)
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("cant start server", "error", err)
		os.Exit(1)
	}
}

func NewServer(port string, h http.Handler) *http.Server {
	slog.Info("starting server", "port", port)
	address := net.JoinHostPort("0.0.0.0", port)

	srv := &http.Server{
		Addr:    address,
		Handler: h,
	}

	return srv
}
