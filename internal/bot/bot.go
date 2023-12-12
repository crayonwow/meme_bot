package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const apiUrl = "https://api.telegram.org/bot{{token}}/{{method_name}}"

type method string

func (m method) String() string { return string(m) }

const methodSendVideo method = "sendVideo"

func newUrl(token string, method method) string {
	s := strings.Replace(apiUrl, "{{token}}", token, 1)
	s = strings.Replace(s, "{{method_name}}", method.String(), 1)
	return s
}

type Bot struct {
	token  string
	client *http.Client
}

func NewBot(token string) *Bot {
	return &Bot{
		token:  token,
		client: &http.Client{},
	}
}

func (b *Bot) UploadVideo(
	ctx context.Context,
	chatID, message string,
	isSilent, hasSpoiler bool,
	videoFile io.Reader,
) error {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		if err := m.WriteField("chat_id", chatID); err != nil {
			w.CloseWithError(err)
			return
		}
		if err := m.WriteField("caption", message); err != nil {
			w.CloseWithError(err)
			return
		}
		if err := m.WriteField("disable_notification", fmt.Sprintf("%t", isSilent)); err != nil {
			w.CloseWithError(err)
			return
		}
		if err := m.WriteField("has_spoiler", fmt.Sprintf("%t", hasSpoiler)); err != nil {
			w.CloseWithError(err)
			return
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

	req, err := http.NewRequest("POST", newUrl(b.token, methodSendVideo), r)
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
