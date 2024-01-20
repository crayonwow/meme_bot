package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"meme_bot/internal/bot"
	"meme_bot/pkg/instagram"
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

	split := strings.Split(message.Message.Text, "\n")

	_url = split[0]
	if len(split) == 2 {
		_message = split[1]
	}
	if l := len(_message); l != 0 {
		shift := 0
		if _isSilent = _message[0] == '.'; _isSilent {
			shift++
		}
		if l > 1 {
			if _hasSpoiler = _message[1] == '.'; _hasSpoiler {
				shift++
			}
		}
		_message = _message[shift:]
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
	rawVideo, err := h.insta.DownloadVideo(ctx, _url)
	if err != nil {
		return fmt.Errorf("cant download video: %w", err)
	}

	const tmpInput = "tmp_in.mp4"
	const tmpOutput = "tmp_out.mp4"

	f, err := os.Create(tmpInput)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, rawVideo)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	stream := ffmpeg.Input(tmpInput)

	err = stream.
		Drawtext(
			"@bruh_memento",
			50,
			50,
			false,
			// ffmpeg.KwArgs{"fontcolor": "red"},
		).
		Output(tmpOutput, ffmpeg.KwArgs{"map": "0:a:?"}).
		ErrorToStdOut().
		OverWriteOutput().
		RunWithResource(0.1, 0.5)
	if err != nil {
		return fmt.Errorf("ffmpeg process: %w", err)
	}

	processedVideo, err := os.Open(tmpOutput)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer processedVideo.Close()

	err = h.b.UploadVideo(ctx, h.chatID, _message, _isSilent, _hasSpoiler, processedVideo)
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
