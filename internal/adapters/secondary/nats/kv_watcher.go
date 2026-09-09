package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go/jetstream"
)

type CameraStateCallback func(tenantID string, cameraID string, state map[string]interface{})

// CameraStateWatcher watches the NATS KV bucket for live status changes.
type CameraStateWatcher struct {
	client *NATSClient
}

// NewCameraStateWatcher creates a new KV state watcher.
func NewCameraStateWatcher(client *NATSClient) *CameraStateWatcher {
	return &CameraStateWatcher{client: client}
}

// SetCameraState updates the live operational state of a camera in the KV bucket.
func (w *CameraStateWatcher) SetCameraState(ctx context.Context, tenantID, cameraID string, state map[string]interface{}) error {
	key := fmt.Sprintf("%s.%s", tenantID, cameraID)
	payload, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal camera state: %w", err)
	}

	_, err = w.client.cameraKV.Put(ctx, key, payload)
	return err
}

// WatchAllCameras streams state changes of all cameras in background.
func (w *CameraStateWatcher) WatchAllCameras(ctx context.Context, cb CameraStateCallback) error {
	watcher, err := w.client.cameraKV.WatchAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to create KV watcher: %w", err)
	}

	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-watcher.Updates():
				if !ok || entry == nil {
					continue
				}

				if entry.Operation() == jetstream.KeyValueDelete || entry.Operation() == jetstream.KeyValuePurge {
					continue
				}

				var state map[string]interface{}
				if err := json.Unmarshal(entry.Value(), &state); err != nil {
					log.Printf("[NATS-KV] Failed to parse camera state for key %s: %v", entry.Key(), err)
					continue
				}

				cb("", entry.Key(), state)
			}
		}
	}()

	return nil
}
