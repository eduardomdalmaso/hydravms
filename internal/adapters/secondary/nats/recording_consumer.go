package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

// RecordingConsumer listens to HYDRA_RECORDINGS JetStream and indexes segments into PostgreSQL.
type RecordingConsumer struct {
	client        *NATSClient
	recordingRepo ports.RecordingRepository
}

// NewRecordingConsumer creates a durable recording consumer.
func NewRecordingConsumer(client *NATSClient, recordingRepo ports.RecordingRepository) *RecordingConsumer {
	return &RecordingConsumer{
		client:        client,
		recordingRepo: recordingRepo,
	}
}

// Start launches a durable JetStream consumer for recording segments.
func (c *RecordingConsumer) Start(ctx context.Context) error {
	if c.client == nil || c.recordingRepo == nil {
		return nil
	}

	consumer, err := c.client.js.CreateOrUpdateConsumer(ctx, "HYDRA_RECORDINGS", jetstream.ConsumerConfig{
		Durable:       "HydraVMS-RecordingIndexer",
		Name:          "HydraVMS-RecordingIndexer",
		Description:   "Durable indexer for camera recording segments",
		FilterSubject: "hydra.v1.*.cameras.*.recordings.segment",
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    5,
		AckWait:       10 * time.Second,
	})
	if err != nil {
		return err
	}

	go func() {
		iter, err := consumer.Messages()
		if err != nil {
			log.Printf("❌ [NATS] Failed to start recording consumer: %v\n", err)
			return
		}
		defer iter.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := iter.Next()
				if err != nil {
					continue
				}

				var payload struct {
					TenantID        string `json:"tenant_id"`
					CameraID        string `json:"camera_id"`
					RecordingMode   string `json:"recording_mode"`
					S3Key           string `json:"s3_key"`
					StartTime       string `json:"start_time"`
					EndTime         string `json:"end_time"`
					DurationSeconds int    `json:"duration_seconds"`
					FileSizeBytes   int64  `json:"file_size_bytes"`
				}

				if err := json.Unmarshal(msg.Data(), &payload); err != nil {
					log.Printf("⚠️ [NATS] Malformed recording segment JSON: %v\n", err)
					msg.Ack()
					continue
				}

				tenantUUID, _ := uuid.Parse(payload.TenantID)
				if tenantUUID == uuid.Nil {
					tenantUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
				}

				startTime, _ := time.Parse(time.RFC3339, payload.StartTime)
				endTime, _ := time.Parse(time.RFC3339, payload.EndTime)
				if startTime.IsZero() {
					startTime = time.Now().Add(-time.Duration(payload.DurationSeconds) * time.Second)
				}
				if endTime.IsZero() {
					endTime = time.Now()
				}

				seg := &domain.RecordingSegment{
					TenantID:        tenantUUID,
					CameraID:        payload.CameraID,
					RecordingMode:   payload.RecordingMode,
					StartTime:       startTime,
					EndTime:         endTime,
					DurationSeconds: payload.DurationSeconds,
					FileSizeBytes:   payload.FileSizeBytes,
					S3Key:           payload.S3Key,
				}

				if err := c.recordingRepo.InsertSegment(context.Background(), seg); err != nil {
					log.Printf("⚠️ [NATS] Failed to index segment for camera %s: %v\n", payload.CameraID, err)
					msg.Nak()
				} else {
					msg.Ack()
				}
			}
		}
	}()

	log.Printf("📦 [NATS] Durable Recording Consumer active on HYDRA_RECORDINGS\n")
	return nil
}
