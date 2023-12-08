package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"meme_bot/internal/bot"
	"meme_bot/pkg/instagram"
	"net/http"
	"net/url"
)

func NewHandler(chatID string, b *bot.Bot) *Handler {
	return &Handler{
		b:      b,
		chatID: chatID,
	}
}

type Handler struct {
	b      *bot.Bot
	chatID string
}

func (h *Handler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read body", http.StatusInternalServerError)
		return
	}
	message := &Message{}
	err = json.Unmarshal(b, message)
	if err != nil {
		http.Error(w, "cant decode payload", http.StatusBadRequest)
		return
	}
	slog.Info("message", "message", string(b))

	_url := message.Message.Text
	_, err = url.Parse(_url)
	if err != nil {
		http.Error(w, "cant parse url", http.StatusBadRequest)
		return
	}
	h.asyncDo(_url)
	slog.Info("success", "url", _url)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) do(ctx context.Context, _url string) error {
	video, err := instagram.DownloadVideo(ctx, _url)
	if err != nil {
		return fmt.Errorf("cant download video: %w", err)
	}

	err = h.b.UploadVideo(ctx, h.chatID, video)
	if err != nil {
		return fmt.Errorf("cant upload video: %w", err)
	}
	return nil
}

func (h *Handler) asyncDo(_url string) {
	go func(__url string) {
		err := h.do(context.Background(), __url)
		if err != nil {
			slog.Error("cant do", "error", err.Error())
		}
	}(_url)
}
