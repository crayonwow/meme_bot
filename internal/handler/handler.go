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
	"slices"
	"strings"
)

func NewHandler(
	client *instagram.Client,
	chatID string,
	b *bot.Bot,
	secretToken string,
	userWhiteList []string,
) *Handler {
	return &Handler{
		b:             b,
		chatID:        chatID,
		insta:         client,
		secretToken:   secretToken,
		userWhiteList: userWhiteList,
	}
}

type Handler struct {
	b             *bot.Bot
	chatID        string
	insta         *instagram.Client
	secretToken   string
	userWhiteList []string
}

func (h *Handler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if token != h.secretToken {
		slog.Error("invalid token", "token", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	slog.Info("message", "message", message, "payload", string(b))
	if !slices.Contains(h.userWhiteList, message.Message.Chat.Username) {
		w.WriteHeader(http.StatusOK)
		return
	}

	var (
		_url, _message         string
		_isSilent, _hasSpoiler bool
	)

	setters := []func(string){
		func(s string) { _url = s },
		func(s string) { _message = s },
		func(s string) { _isSilent = s == "1" },
		func(s string) { _hasSpoiler = s == "1" },
	}

	for i, s := range strings.Split(message.Message.Text, "\n") {
		setters[i](s)
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
