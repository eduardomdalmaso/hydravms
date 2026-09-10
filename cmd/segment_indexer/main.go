package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if len(os.Args) < 3 {
		return
	}
	camID := os.Args[1]
	segPath := os.Args[2]
	var durationSec int = 30
	if len(os.Args) >= 4 {
		if d, err := strconv.ParseFloat(os.Args[3], 64); err == nil && d > 0 {
			durationSec = int(d)
		}
	}

	info, err := os.Stat(segPath)
	var fileSizeBytes int64 = 1048576
	if err == nil {
		fileSizeBytes = info.Size()
	}

	base := filepath.Base(segPath)
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})_(\d{2})-(\d{2})-(\d{2})`)
	var startTime time.Time
	if match := re.FindStringSubmatch(base); len(match) == 4 {
		loc := time.Local
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s:%s:%s", match[1], match[2], match[3], match[4]), loc)
	}
	if startTime.IsZero() {
		startTime = time.Now().Add(-time.Duration(durationSec) * time.Second)
	}
	endTime := startTime.Add(time.Duration(durationSec) * time.Second)

	nc, err := nats.Connect("nats://localhost:4222", nats.Timeout(1*time.Second))
	if err != nil {
		return
	}
	defer nc.Close()

	payload, _ := json.Marshal(map[string]interface{}{
		"tenant_id":        "00000000-0000-0000-0000-000000000001",
		"camera_id":        camID,
		"recording_mode":   "motion",
		"s3_key":           segPath,
		"start_time":       startTime.Format(time.RFC3339),
		"end_time":         endTime.Format(time.RFC3339),
		"duration_seconds": durationSec,
		"file_size_bytes":  fileSizeBytes,
	})

	subject := fmt.Sprintf("hydra.v1.00000000-0000-0000-0000-000000000001.cameras.%s.recordings.segment", camID)
	_ = nc.Publish(subject, payload)
	_ = nc.Flush()
}
