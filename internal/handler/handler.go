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
	"strings"
)

func NewHandler(client *instagram.Client, chatID string, b *bot.Bot) *Handler {
	return &Handler{
		b:      b,
		chatID: chatID,
		insta:  client,
	}
}

type Handler struct {
	b      *bot.Bot
	chatID string
	insta  *instagram.Client
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
		slog.Error("cant decode payload", "error", err.Error(), "payload", string(b))
		w.WriteHeader(http.StatusOK)
		return
	}
	slog.Info("message", "message", string(b))

	text := message.Message.Text
	splits := strings.Split(text, "\n")
	if len(splits) != 4 {
		slog.Error("cant parse message", "splits", splits)
		w.WriteHeader(http.StatusOK)
		return
	}

	_url := splits[0]
	_message := splits[1]
	_isSilent := splits[2] == "1"
	_hasSpoiler := splits[3] == "1"

	_, err = url.Parse(_url)
	if err != nil {
		slog.Error("cant parse url", "url", _url)
		w.WriteHeader(http.StatusOK)
		return
	}
	h.asyncDo(_url, _message, _isSilent, _hasSpoiler)
	slog.Info("success", "url", _url)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) do(
	ctx context.Context,
	_url, _message string,
	_isSilent, _hasSpoiler bool,
) error {
	video, err := h.insta.DownloadVideo(ctx, _url)
	if err != nil {
		return fmt.Errorf("cant download video: %w", err)
	}

	err = h.b.UploadVideo(ctx, h.chatID, _message, _isSilent, _hasSpoiler, video)
	if err != nil {
		return fmt.Errorf("cant upload video: %w", err)
	}
	return nil
}

func (h *Handler) asyncDo(_url, _mesage string, _isSilent, _hasSpoiler bool) {
	go func(_url, _mesage string, _isSilent, _hasSpoiler bool) {
		err := h.do(context.Background(), _url, _mesage, _isSilent, _hasSpoiler)
		if err != nil {
			slog.Error("cant do", "error", err.Error())
		}
	}(_url, _mesage, _isSilent, _hasSpoiler)
}
