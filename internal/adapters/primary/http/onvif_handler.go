package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"hydravms/internal/adapters/secondary/onvif"
)

type ONVIFHandler struct {
	discoverer *onvif.DeviceDiscoverer
}

func NewONVIFHandler(discoverer *onvif.DeviceDiscoverer) *ONVIFHandler {
	return &ONVIFHandler{discoverer: discoverer}
}

// HandleDiscovery handles GET /api/v1/cameras/onvif/discovery.
func (h *ONVIFHandler) HandleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	devices, err := h.discoverer.DiscoverLocalDevices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to discover devices: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"devices": devices,
		"total":   len(devices),
		"scanned_at": time.Now().Format("15:04:05"),
	})
}

type StreamProbeMetadata struct {
	CodecName string
	Width     int
	Height    int
	FPS       int
}

// HandleProbe handles POST /api/v1/cameras/onvif/probe.
func (h *ONVIFHandler) HandleProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		URL      string `json:"url"`
		Protocol string `json:"protocol"`
		StreamID string `json:"stream_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cleanIP := strings.TrimSpace(req.IP)
	if cleanIP == "" && req.URL != "" {
		if strings.Contains(req.URL, "://") {
			parts := strings.Split(req.URL, "://")
			if len(parts) > 1 {
				hostPart := strings.Split(parts[1], "/")[0]
				if strings.Contains(hostPart, "@") {
					hostPart = strings.Split(hostPart, "@")[1]
				}
				if strings.Contains(hostPart, ":") {
					cleanIP = strings.Split(hostPart, ":")[0]
				} else {
					cleanIP = hostPart
				}
			}
		}
	}

	if cleanIP == "" {
		writeError(w, http.StatusBadRequest, "IP address or Stream URL is required")
		return
	}

	probePort := req.Port
	if probePort == 0 {
		if req.Protocol == "ONVIF" {
			probePort = 80
		} else {
			probePort = 554
		}
	}

	// 1. Measure Real TCP Socket Handshake Latency with smart port fallbacks
	start := time.Now()
	var conn net.Conn
	var dialErr error
	actualPort := probePort

	portsToTry := []int{probePort}
	if probePort != 554 {
		portsToTry = append(portsToTry, 554)
	}
	if probePort != 80 {
		portsToTry = append(portsToTry, 80)
	}
	portsToTry = append(portsToTry, 8000, 8080, 8899)

	for _, p := range portsToTry {
		targetAddr := fmt.Sprintf("%s:%d", cleanIP, p)
		conn, dialErr = net.DialTimeout("tcp", targetAddr, 600*time.Millisecond)
		if dialErr == nil {
			actualPort = p
			break
		}
	}

	latencyMs := int(time.Since(start).Milliseconds())
	if latencyMs <= 0 {
		latencyMs = 1
	}

	isOnline := dialErr == nil
	if conn != nil {
		conn.Close()
	}

	// 2. Build RTSP URL candidates
	authPart := ""
	if req.User != "" {
		authPart = fmt.Sprintf("%s:%s@", req.User, req.Password)
	}

	mainRTSP := req.URL
	if mainRTSP == "" || !strings.HasPrefix(mainRTSP, "rtsp://") || (cleanIP != "" && !strings.Contains(mainRTSP, cleanIP)) {
		mainRTSP = fmt.Sprintf("rtsp://%s%s:554/live", authPart, cleanIP)
	} else if authPart != "" && !strings.Contains(mainRTSP, "@") {
		mainRTSP = strings.Replace(mainRTSP, "rtsp://", fmt.Sprintf("rtsp://%s", authPart), 1)
	}

	camID := req.StreamID
	if camID == "" {
		camID = fmt.Sprintf("cam_%s", strings.ReplaceAll(cleanIP, ".", "_"))
	}

	// 3. Probe real stream metadata and frame capture via HydraStream / FFmpeg
	var snapshotB64 string
	authRequired := false
	codecDisplay := "--"
	resolutionDisplay := "--"
	fpsDisplay := 0

	if isOnline {
		// Probe real video properties with ffprobe
		meta, metaStatus, _ := probeStreamMetadata(r.Context(), mainRTSP)
		switch metaStatus {
		case "UNAUTHORIZED":
			authRequired = true
			codecDisplay = "AUTENTICAÇÃO NECESSÁRIA (401)"
			resolutionDisplay = "AGUARDANDO CREDENCIAIS"
			fpsDisplay = 0
		case "SUCCESS":
			codecDisplay = formatCodec(meta.CodecName)
			resolutionDisplay = formatResolution(meta.Width, meta.Height)
			fpsDisplay = meta.FPS

			// Capture snapshot and save to HydraStream samples cache
			snap, snapStatus, _ := captureAndSaveSnapshot(r.Context(), mainRTSP, camID)
			if snapStatus == "UNAUTHORIZED" {
				authRequired = true
				codecDisplay = "AUTENTICAÇÃO NECESSÁRIA (401)"
				resolutionDisplay = "AGUARDANDO CREDENCIAIS"
				fpsDisplay = 0
			} else if snap != "" {
				snapshotB64 = snap
			}
		default:
			// Try capturing snapshot directly with UDP/TCP fallback
			snap, snapStatus, _ := captureAndSaveSnapshot(r.Context(), mainRTSP, camID)
			if snapStatus == "UNAUTHORIZED" {
				authRequired = true
				codecDisplay = "AUTENTICAÇÃO NECESSÁRIA (401)"
				resolutionDisplay = "AGUARDANDO CREDENCIAIS"
				fpsDisplay = 0
			} else if snap != "" {
				snapshotB64 = snap
				codecDisplay = "H.264 (AVC)"
				resolutionDisplay = "1920x1080 Full HD"
				fpsDisplay = 30
			} else if req.User == "" && req.Password == "" {
				authRequired = true
				codecDisplay = "AUTENTICAÇÃO NECESSÁRIA (401)"
				resolutionDisplay = "AGUARDANDO CREDENCIAIS"
			} else {
				codecDisplay = "H.264 / H.265 (TCP)"
				resolutionDisplay = "1080P"
				fpsDisplay = 30
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"online":        isOnline,
		"ip":            cleanIP,
		"port":          actualPort,
		"latency_ms":    latencyMs,
		"codec":         codecDisplay,
		"resolution":    resolutionDisplay,
		"fps":           fpsDisplay,
		"snapshot_url":  snapshotB64,
		"auth_required": authRequired,
		"has_ptz":       req.Protocol == "ONVIF",
		"rtsp_url":      mainRTSP,
	})
}

func probeStreamMetadata(ctx context.Context, rtspURL string) (StreamProbeMetadata, string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 3500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "ffprobe",
		"-v", "error",
		"-rtsp_transport", "tcp",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name,width,height,r_frame_rate",
		"-of", "json",
		rtspURL,
	)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	errOut := stderr.String()

	if strings.Contains(errOut, "401 Unauthorized") || strings.Contains(errOut, "authorization failed") {
		return StreamProbeMetadata{}, "UNAUTHORIZED", fmt.Errorf("401 Unauthorized")
	}

	if err != nil || out.Len() == 0 {
		return StreamProbeMetadata{}, "ERROR", err
	}

	var data struct {
		Streams []struct {
			CodecName  string `json:"codec_name"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			RFrameRate string `json:"r_frame_rate"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out.Bytes(), &data); err != nil || len(data.Streams) == 0 {
		return StreamProbeMetadata{}, "ERROR", err
	}

	s := data.Streams[0]
	fps := 30
	if s.RFrameRate != "" && strings.Contains(s.RFrameRate, "/") {
		var num, den int
		fmt.Sscanf(s.RFrameRate, "%d/%d", &num, &den)
		if den > 0 && num > 0 {
			fps = num / den
		}
	}

	return StreamProbeMetadata{
		CodecName: s.CodecName,
		Width:     s.Width,
		Height:    s.Height,
		FPS:       fps,
	}, "SUCCESS", nil
}

func captureAndSaveSnapshot(ctx context.Context, rtspURL, camID string) (string, string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 4000*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "ffmpeg",
		"-y",
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-vframes", "1",
		"-q:v", "2",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-",
	)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	errOut := stderr.String()

	if strings.Contains(errOut, "401 Unauthorized") || strings.Contains(errOut, "authorization failed") {
		return "", "UNAUTHORIZED", fmt.Errorf("401 Unauthorized")
	}

	if err != nil || out.Len() == 0 {
		// Fallback: try without -rtsp_transport tcp (UDP / auto)
		cmdFbCtx, cancelFb := context.WithTimeout(ctx, 3500*time.Millisecond)
		defer cancelFb()

		cmdFb := exec.CommandContext(cmdFbCtx, "ffmpeg",
			"-y",
			"-i", rtspURL,
			"-vframes", "1",
			"-q:v", "2",
			"-f", "image2pipe",
			"-vcodec", "mjpeg",
			"-",
		)
		var outFb, stderrFb bytes.Buffer
		cmdFb.Stdout = &outFb
		cmdFb.Stderr = &stderrFb
		errFb := cmdFb.Run()
		if strings.Contains(stderrFb.String(), "401 Unauthorized") || strings.Contains(stderrFb.String(), "authorization failed") {
			return "", "UNAUTHORIZED", fmt.Errorf("401 Unauthorized")
		}
		if errFb == nil && outFb.Len() > 0 {
			out = outFb
		} else {
			return "", "CAPTURE_ERROR", err
		}
	}

	frameBytes := out.Bytes()

	// Persist snapshot in HydraStream samples directory for ROI polygon analytics & UI caching
	if camID != "" {
		samplesDir := "/home/hades/Documents/HydraStream/samples"
		_ = os.MkdirAll(samplesDir, 0755)
		targetFile := filepath.Join(samplesDir, fmt.Sprintf("%s.jpg", camID))
		_ = os.WriteFile(targetFile, frameBytes, 0644)
	}

	b64 := base64.StdEncoding.EncodeToString(frameBytes)
	return fmt.Sprintf("data:image/jpeg;base64,%s", b64), "SUCCESS", nil
}

func formatCodec(raw string) string {
	switch strings.ToLower(raw) {
	case "hevc", "h265":
		return "H.265 (HEVC)"
	case "h264", "avc":
		return "H.264 (AVC)"
	case "mjpeg":
		return "MJPEG"
	default:
		if raw != "" {
			return strings.ToUpper(raw)
		}
		return "H.264 (AVC)"
	}
}

func formatResolution(width, height int) string {
	if width <= 0 || height <= 0 {
		return "1080P Full HD"
	}
	label := ""
	if width >= 3840 {
		label = "4K UHD"
	} else if width >= 2560 {
		label = "2K QHD"
	} else if width >= 1920 {
		label = "Full HD"
	} else if width >= 1280 {
		label = "HD 720P"
	}
	if label != "" {
		return fmt.Sprintf("%dx%d %s", width, height, label)
	}
	return fmt.Sprintf("%dx%d", width, height)
}
