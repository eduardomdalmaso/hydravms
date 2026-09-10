package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"hydravms/internal/domain"
)

// EventPublisher publishes AI events and camera telemetry to JetStream.
type EventPublisher struct {
	client *NATSClient
}

// NewEventPublisher creates a new JetStream publisher.
func NewEventPublisher(client *NATSClient) *EventPublisher {
	return &EventPublisher{client: client}
}

// PublishAIEvent sends an event to NATS with deduplication ID.
func (p *EventPublisher) PublishAIEvent(ctx context.Context, event *domain.Event) error {
	subject := fmt.Sprintf("hydra.v1.%s.cameras.%s.events", event.TenantID.String(), event.CameraID)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    payload,
		Header:  nats.Header{},
	}
	// JetStream deduplication header
	msg.Header.Set(jetstream.MsgIDHeader, fmt.Sprintf("evt_%s", event.ID.String()))

	_, err = p.client.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish AI event to %s: %w", subject, err)
	}

	return nil
}

// PublishCloudEvent sends a CNCF CloudEvents standard envelope to NATS.
func (p *EventPublisher) PublishCloudEvent(ctx context.Context, event *domain.CloudEvent) error {
	subject := fmt.Sprintf("hydra.v1.%s.events.%s", event.TenantID.String(), event.Type)
	if event.Subject != "" {
		subject = fmt.Sprintf("hydra.v1.%s.cameras.%s.events", event.TenantID.String(), event.Subject)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal cloud event: %w", err)
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    payload,
		Header:  nats.Header{},
	}
	msg.Header.Set(jetstream.MsgIDHeader, event.ID)

	if p.client.js != nil {
		if _, err := p.client.js.PublishMsg(ctx, msg); err == nil {
			return nil
		}
	}
	return p.client.nc.Publish(subject, payload)
}

// PublishTelemetry sends high-frequency camera FPS and bitrate metrics.
func (p *EventPublisher) PublishTelemetry(ctx context.Context, tenantID string, cameraID string, telemetry map[string]interface{}) error {
	subject := fmt.Sprintf("hydra.v1.%s.cameras.%s.telemetry", tenantID, cameraID)

	payload, err := json.Marshal(telemetry)
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry: %w", err)
	}

	return p.client.nc.Publish(subject, payload)
}

