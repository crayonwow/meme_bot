package main

import (
	"flag"
	"log/slog"
	"meme_bot/internal/bot"
	"meme_bot/internal/handler"
	"net"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "8000", "port number for http listener")
	token := flag.String("token", "", "bot token")
	chatID := flag.String("chat_id", "", "")
	flag.Parse()

	if *token == "" {
		slog.Error("token is empty")
		os.Exit(1)
	}
	if *chatID == "" {
		slog.Error("chatID is empty")
		os.Exit(1)
	}

	b := bot.NewBot(*token)
	h := handler.NewHandler(*chatID, b)

	router := http.NewServeMux()
	router.HandleFunc("/webhook", h.HandleMessage)

	srv := NewServer(*port, router)
	err := srv.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}
}

func NewServer(port string, h http.Handler) *http.Server {
	address := net.JoinHostPort("0.0.0.0", port)

	srv := &http.Server{
		Addr:    address,
		Handler: h,
	}

	return srv
}
