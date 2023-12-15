package main

import (
	"log/slog"
	"meme_bot/internal/bot"
	"meme_bot/internal/handler"
	"meme_bot/pkg"
	"meme_bot/pkg/instagram"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	token := os.Getenv("TOKEN")
	chatID := os.Getenv("CHAT_ID")
	secretToken := os.Getenv("SECRET_TOKEN")
	userWhiteListRaw := os.Getenv("USER_WHITE_LIST")

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

	if secretToken == "" {
		slog.Error("SECRET_TOKEN is empty")
		os.Exit(1)
	}

	if userWhiteListRaw == "" || len(strings.Split(userWhiteListRaw, ",")) == 0 {
		slog.Error("WHITE_LIST is empty")
		os.Exit(1)
	}

	client := pkg.NewHttpClient()
	insta := instagram.NewClient(client)

	b := bot.NewBot(client, token, time.Minute)
	h := handler.NewHandler(insta, chatID, b, secretToken, strings.Split(userWhiteListRaw, ","))

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
