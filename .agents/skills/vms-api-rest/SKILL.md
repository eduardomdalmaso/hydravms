---
name: vms-api-rest
description: REST API conventions, versioning (/api/v1), Swagger/OpenAPI 3.0 specs, JWT authentication headers, RFC 7807 error responses, and WebSocket gateways for HydraVMS.
---

# 🔌 REST API & Swagger / OpenAPI Standards Skill (HydraVMS)

This skill defines the REST API standards, endpoints catalog, Swagger OpenAPI documentation, and real-time WebSocket protocol for HydraVMS.

---

## 1. REST API Conventions & Versioning

- **Base URL:** `/api/v1`
- **Authentication:** `Authorization: Bearer <jwt_token>` header on every request.
- **Content-Type:** `application/json; charset=utf-8` (or `multipart/form-data` for uploads).
- **Standard Error Response (RFC 7807 Problem Details):**
```json
{
  "type": "https://hydravms.io/errors/not-found",
  "title": "Camera Not Found",
  "status": 404,
  "detail": "Camera 'cam_portaria_01' does not exist for tenant 'tenant_123'",
  "instance": "/api/v1/cameras/cam_portaria_01",
  "timestamp": "2026-09-03T19:30:00Z"
}
```

---

## 2. Endpoints Catalog & Swagger UI

The interactive Swagger UI is available at `GET /swagger/index.html` generated from `api/openapi.yaml`.

### 🔑 Authentication & Tenants
- `POST /api/v1/auth/login` - Authenticate user & issue JWT.
- `POST /api/v1/auth/refresh` - Refresh active JWT session.
- `GET /api/v1/auth/me` - Get current user profile and permissions.
- `GET /api/v1/tenants` - List tenants (Master Admin).
- `POST /api/v1/tenants` - Create new tenant with camera quotas.

### 🎥 Cameras & Streaming
- `GET /api/v1/cameras` - List cameras for active tenant (status, FPS, bitrate).
- `POST /api/v1/cameras` - Register new camera with RTSP URL.
- `GET /api/v1/cameras/:id` - Get camera details and stream health.
- `PUT /api/v1/cameras/:id` - Update camera properties and recording profile.
- `DELETE /api/v1/cameras/:id` - Remove camera and cascade active recordings.
- `GET /api/v1/cameras/:id/snapshot` - Get real-time snapshot JPEG from stream.
- `POST /api/v1/cameras/:id/whep` - WebRTC WHEP endpoint for live streaming playback.

### 🎯 AI Rules, Zones & SAHI
- `GET /api/v1/cameras/:id/zones` - List detection zones (polygons) for camera.
- `POST /api/v1/cameras/:id/zones` - Create new intrusion/tripwire zone.
- `GET /api/v1/rules` - List AI alarm rules (classes, confidence, SAHI config).
- `POST /api/v1/rules` - Create AI detection rule with trigger actions.
- `PUT /api/v1/rules/:id` - Update rule sensitivity, cooldown or SAHI slicing parameters.
- `DELETE /api/v1/rules/:id` - Remove AI rule.

### 🚨 AI Events & Incidents (Gold Layer)
- `GET /api/v1/events` - Query events with filters (`camera_id`, `event_type`, `severity`, `status`, `from`, `to`).
- `GET /api/v1/events/:id` - Get event details with bbox, confidence and metadata.
- `POST /api/v1/events/:id/acknowledge` - Operator acknowledge alarm.
- `POST /api/v1/events/:id/resolve` - Resolve event with resolution notes.
- `POST /api/v1/events/:id/pin` - Flag event as legal proof / pinned evidence (exempt from purge).
- `GET /api/v1/events/:id/clip-url` - Generate temporary MinIO Presigned URL for video playback.

### 📼 Recordings & Playback
- `GET /api/v1/recordings/timeline` - Get timeline blocks for camera and date range.
- `GET /api/v1/recordings/:id/stream-url` - Generate Presigned URL for fMP4 segment playback.

### 📤 Video Exports
- `POST /api/v1/exports` - Request video clip export (Zero-Copy or NVENC Watermark).
- `GET /api/v1/exports` - List export jobs and progress status.
- `GET /api/v1/exports/:id/download-url` - Get Presigned download URL for exported MP4.

### 💾 Storage & Hardware Telemetry
- `GET /api/v1/storage/pools` - List HD volumes, used/free bytes, and watermark status.
- `POST /api/v1/storage/pools` - Register new storage drive mount path into MinIO pool.
- `POST /api/v1/storage/purge` - Manually trigger storage purge.
- `GET /api/v1/telemetry/system` - Get real-time CPU, RAM, GPU (NVML), and StorageGuard metrics.
- `GET /api/v1/cluster/nodes` - List all connected cluster nodes (Ingest, GPU Workers, MinIO).

---

## 3. Real-Time WebSocket Gateway

- **Endpoint:** `GET /ws/v1/live-feed?token=<jwt_token>`
- **Payload Events:**
  - `event.alarm.new`: Instant push of Gold Layer AI event with crop thumbnail.
  - `event.alarm.ack`: Broadcast when an operator acknowledges an alert.
  - `telemetry.tick`: 1Hz hardware and stream status updates.
  - `camera.status.change`: Online / Offline status notification.
