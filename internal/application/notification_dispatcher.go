package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"hydravms/internal/domain"
)

// NotificationDispatcher dispatches real-time security alerts to external webhooks and Telegram bots.
type NotificationDispatcher struct {
	client         *http.Client
	webhookURL     string
	telegramToken  string
	telegramChatID string
}

// NewNotificationDispatcher creates an instance reading configuration from environment variables.
func NewNotificationDispatcher() *NotificationDispatcher {
	return &NotificationDispatcher{
		client: &http.Client{
			Timeout: 4 * time.Second,
		},
		webhookURL:     strings.TrimSpace(os.Getenv("HYDRA_ALARM_WEBHOOK_URL")),
		telegramToken:  strings.TrimSpace(os.Getenv("HYDRA_TELEGRAM_BOT_TOKEN")),
		telegramChatID: strings.TrimSpace(os.Getenv("HYDRA_TELEGRAM_CHAT_ID")),
	}
}

// AlarmPayload represents the standardized JSON structure delivered to HTTP webhooks.
type AlarmPayload struct {
	EventID     string    `json:"event_id"`
	CameraID    string    `json:"camera_id"`
	CameraName  string    `json:"camera_name"`
	EventType   string    `json:"event_type"`
	Severity    string    `json:"severity"`
	Confidence  float64   `json:"confidence"`
	TriggeredAt time.Time `json:"triggered_at"`
	SnapshotURL string    `json:"snapshot_url,omitempty"`
}

// Dispatch sends notifications asynchronously to all configured external channels.
func (d *NotificationDispatcher) Dispatch(evt *domain.Event) {
	if evt == nil {
		return
	}

	payload := AlarmPayload{
		EventID:     evt.ID.String(),
		CameraID:    evt.CameraID,
		CameraName:  evt.CameraName,
		EventType:   evt.EventType,
		Severity:    string(evt.Severity),
		Confidence:  evt.Confidence,
		TriggeredAt: evt.TriggeredAt,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if d.webhookURL != "" {
			if err := d.sendWebhook(ctx, payload); err != nil {
				log.Printf("⚠️ [Notifier] Webhook dispatch failed: %v", err)
			}
		}

		if d.telegramToken != "" && d.telegramChatID != "" {
			if err := d.sendTelegram(ctx, payload); err != nil {
				log.Printf("⚠️ [Notifier] Telegram dispatch failed: %v", err)
			}
		}
	}()
}

func (d *NotificationDispatcher) sendWebhook(ctx context.Context, payload AlarmPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.webhookURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "HydraVMS-Notifier/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook responded with status: %d", resp.StatusCode)
	}
	return nil
}

func (d *NotificationDispatcher) sendTelegram(ctx context.Context, payload AlarmPayload) error {
	text := fmt.Sprintf(
		"🚨 *[ALERTA HYDRA]*\n*Câmera:* %s (`%s`)\n*Tipo:* %s\n*Severidade:* %s\n*Horário:* %s",
		payload.CameraName,
		payload.CameraID,
		payload.EventType,
		strings.ToUpper(payload.Severity),
		payload.TriggeredAt.Format("02/01/2006 15:04:05"),
	)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", d.telegramToken)
	form := url.Values{}
	form.Set("chat_id", d.telegramChatID)
	form.Set("text", text)
	form.Set("parse_mode", "Markdown")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("telegram request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram responded with status: %d", resp.StatusCode)
	}
	return nil
}
