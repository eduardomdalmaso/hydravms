package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATSClient manages NATS connection, JetStream contexts, and KV buckets.
type NATSClient struct {
	nc         *nats.Conn
	js         jetstream.JetStream
	cameraKV   jetstream.KeyValue
}

// Config holds NATS connection and stream parameters.
type Config struct {
	URL            string
	MaxReconnects  int
	ReconnectWait  time.Duration
}

// DefaultConfig returns production-ready default parameters.
func DefaultConfig() Config {
	return Config{
		URL:           "nats://localhost:4222",
		MaxReconnects: -1, // Infinite reconnect
		ReconnectWait: 2 * time.Second,
	}
}

// NewNATSClient connects to NATS and ensures all Hydra JetStream streams and KV buckets exist.
func NewNATSClient(ctx context.Context, cfg Config) (*NATSClient, error) {
	opts := []nats.Option{
		nats.Name("HydraVMS-ControlPlane"),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.PingInterval(20 * time.Second),
		nats.MaxPingsOutstanding(3),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS at %s: %w", cfg.URL, err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	}

	client := &NATSClient{nc: nc, js: js}

	if err := client.initStreamsAndBuckets(ctx); err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to initialize streams: %w", err)
	}

	return client, nil
}

func (c *NATSClient) initStreamsAndBuckets(ctx context.Context) error {
	// 1. Persistent Events Stream (File Storage, 7-Day Retention)
	eventsStreamCfg := jetstream.StreamConfig{
		Name:        "HYDRA_EVENTS",
		Description: "Gold Layer AI Events, Alarms and Forensic Detections",
		Subjects:    []string{"hydra.v1.*.cameras.*.events", "hydra.v1.*.alarms.>"},
		Retention:   jetstream.LimitsPolicy,
		Storage:     jetstream.FileStorage,
		MaxAge:      7 * 24 * time.Hour,
		Duplicates:  5 * time.Minute,
	}
	if _, err := c.js.CreateOrUpdateStream(ctx, eventsStreamCfg); err != nil {
		return fmt.Errorf("failed to create HYDRA_EVENTS stream: %w", err)
	}

	// 2. High-Throughput Telemetry Stream (Memory Storage, 10-Minute Rolling)
	telemetryStreamCfg := jetstream.StreamConfig{
		Name:        "HYDRA_TELEMETRY",
		Description: "Real-time Camera FPS, Bitrate and Node Telemetry",
		Subjects:    []string{"hydra.v1.*.cameras.*.telemetry", "hydra.v1.cluster.nodes.*.heartbeat"},
		Retention:   jetstream.LimitsPolicy,
		Storage:     jetstream.MemoryStorage,
		MaxAge:      10 * time.Minute,
	}
	if _, err := c.js.CreateOrUpdateStream(ctx, telemetryStreamCfg); err != nil {
		return fmt.Errorf("failed to create HYDRA_TELEMETRY stream: %w", err)
	}

	// 3. Persistent Recordings Stream (File Storage, 7-Day Retention)
	recordingsStreamCfg := jetstream.StreamConfig{
		Name:        "HYDRA_RECORDINGS",
		Description: "Camera Video Recording Segments and Storage Metadata",
		Subjects:    []string{"hydra.v1.*.cameras.*.recordings.segment", "hydra.recordings.>"},
		Retention:   jetstream.LimitsPolicy,
		Storage:     jetstream.FileStorage,
		MaxAge:      7 * 24 * time.Hour,
	}
	if _, err := c.js.CreateOrUpdateStream(ctx, recordingsStreamCfg); err != nil {
		return fmt.Errorf("failed to create HYDRA_RECORDINGS stream: %w", err)
	}

	// 4. Camera States KV Bucket for Live Status Cache
	kvCfg := jetstream.KeyValueConfig{
		Bucket:      "HYDRA_CAMERA_STATES",
		Description: "Real-time operational states and health of cameras",
		History:     5,
		TTL:         24 * time.Hour,
		Storage:     jetstream.MemoryStorage,
	}
	kv, err := c.js.CreateOrUpdateKeyValue(ctx, kvCfg)
	if err != nil {
		return fmt.Errorf("failed to create HYDRA_CAMERA_STATES KV bucket: %w", err)
	}
	c.cameraKV = kv

	return nil
}

// Conn returns the raw NATS connection.
func (c *NATSClient) Conn() *nats.Conn {
	return c.nc
}

// JetStream returns the JetStream interface.
func (c *NATSClient) JetStream() jetstream.JetStream {
	return c.js
}

// CameraKV returns the KeyValue bucket for cameras.
func (c *NATSClient) CameraKV() jetstream.KeyValue {
	return c.cameraKV
}

// Close gracefully closes the NATS connection.
func (c *NATSClient) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}
