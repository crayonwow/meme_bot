package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const apiUrl = "https://api.telegram.org/bot{{token}}/{{method_name}}"

const methodSendVideo = "sendVideo"

func newUrl(token string, method string) string {
	s := strings.Replace(apiUrl, "{{token}}", token, 1)
	s = strings.Replace(s, "{{method_name}}", method, 1)
	return s
}

type Bot struct {
	token   string
	client  *http.Client
	limiter *rate.Limiter
}

func NewBot(cli *http.Client, token string, notifyEvery time.Duration) *Bot {
	return &Bot{
		token:   token,
		client:  cli,
		limiter: rate.NewLimiter(rate.Every(1*notifyEvery/1), 1),
	}
}

func (b *Bot) UploadVideo(
	ctx context.Context,
	chatID, message string,
	isSilent, hasSpoiler bool,
	videoFile io.Reader,
) error {
	return b.uploadVideo(ctx, chatID, message, b.isSilent(isSilent), hasSpoiler, videoFile)
}

func (b *Bot) isSilent(isSilent bool) bool {
	if !isSilent {
		isSilent = !b.limiter.Allow()
	}
	return isSilent
}

func (b *Bot) uploadVideo(
	ctx context.Context,
	chatID, message string,
	isSilent, hasSpoiler bool,
	videoFile io.Reader,
) error {
	slog.Info("uploading video",
		"chatID", chatID,
		"message", message,
		"isSilent", isSilent,
		"hasSpoiler", hasSpoiler,
	)

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		for k, v := range map[string]string{
			"chat_id":              chatID,
			"caption":              message,
			"disable_notification": fmt.Sprintf("%t", isSilent),
			"has_spoiler":          fmt.Sprintf("%t", hasSpoiler),
		} {
			if err := m.WriteField(k, v); err != nil {
				w.CloseWithError(err)
				return
			}
		}
		part, err := m.CreateFormFile("video", "video.mp4")
		if err != nil {
			w.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, videoFile); err != nil {
			w.CloseWithError(err)
			return
		}

		if closer, ok := videoFile.(io.ReadCloser); ok {
			if err = closer.Close(); err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()

	req, err := http.NewRequest(http.MethodPost, newUrl(b.token, methodSendVideo), r)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Ok {
		return fmt.Errorf("api error: %s", apiResp.Description)
	}

	return nil
}
