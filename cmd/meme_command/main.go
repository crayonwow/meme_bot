package main

import (
	"context"
	"flag"
	"log/slog"
	"meme_bot/internal/bot"
	"meme_bot/pkg/instagram"
	"os"
)

func main() {
	_token := flag.String("token", "", "telegram bot token")
	_chatID := flag.String("chat_id", "", "telegram chat id")
	_url := flag.String("url", "", "instagram url")
	flag.Parse()

	token := *_token
	chatID := *_chatID
	url := *_url

	slog.Info("starting app", "token", token, "chatID", chatID)

	if token == "" {
		slog.Error("token is empty")
		os.Exit(1)
	}
	if chatID == "" {
		slog.Error("chatID is empty")
		os.Exit(1)
	}

	b := bot.NewBot(token)

	ctx := context.Background()
	video, err := instagram.DownloadVideo(ctx, url)
	if err != nil {
		slog.Error("cant download video", "err", err.Error())
		os.Exit(1)
	}

	err = b.UploadVideo(ctx, chatID, video)
	if err != nil {
		slog.Error("cant upload video", "err", err.Error())
		os.Exit(1)
	}
}
