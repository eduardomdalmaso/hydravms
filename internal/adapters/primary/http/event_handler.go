package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/adapters/primary/http/middleware"
	"hydravms/internal/adapters/primary/ws"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type EventHandler struct {
	eventRepo ports.EventRepository
	hub       *ws.Hub
}

func NewEventHandler(eventRepo ports.EventRepository, hub *ws.Hub) *EventHandler {
	return &EventHandler{eventRepo: eventRepo, hub: hub}
}

func (h *EventHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	switch r.Method {
	case http.MethodGet:
		camID := r.URL.Query().Get("camera_id")
		events, err := h.eventRepo.List(r.Context(), tenantID, camID, 100)
		if err != nil {
			http.Error(w, `{"error":"failed to list events"}`, http.StatusInternalServerError)
			return
		}
		if events == nil {
			events = []*domain.Event{}
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"events": events,
			"count":  len(events),
		})

	case http.MethodPost:
		var req struct {
			CameraID    string              `json:"camera_id"`
			EventType   string              `json:"event_type"`
			Severity    domain.EventSeverity `json:"severity"`
			ObjectClass string              `json:"object_class"`
			Confidence  float64             `json:"confidence"`
			BBox        domain.BoundingBox  `json:"bbox_normalized"`
			SnapshotURL string              `json:"snapshot_url"`
			Notes       string              `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.Severity == "" {
			req.Severity = domain.SeverityMedium
		}
		if req.EventType == "" {
			req.EventType = "ai.detection.person"
		}

		evt := &domain.Event{
			ID:             uuid.New(),
			TenantID:       tenantID,
			CameraID:       req.CameraID,
			EventType:      req.EventType,
			Severity:       req.Severity,
			Status:         domain.EventStatusNew,
			TriggeredAt:    time.Now().UTC(),
			ObjectClass:    req.ObjectClass,
			Confidence:     req.Confidence,
			BBoxNormalized: req.BBox,
			SnapshotS3Key:  req.SnapshotURL,
			Notes:          req.Notes,
			CreatedAt:      time.Now().UTC(),
		}

		if err := h.eventRepo.Create(r.Context(), evt); err != nil {
			http.Error(w, `{"error":"failed to save event"}`, http.StatusInternalServerError)
			return
		}

		if h.hub != nil {
			cloudEvt := domain.NewCloudEvent(
				tenantID,
				req.EventType,
				domain.CategoryAIEvent,
				req.Severity,
				req.CameraID,
				map[string]interface{}{
					"camera_id":    req.CameraID,
					"camera_name":  "C113",
					"class_name":   req.ObjectClass,
					"confidence":   req.Confidence,
					"bbox":         req.BBox,
					"snapshot_url": req.SnapshotURL,
					"notes":        req.Notes,
				},
			)
			if data, err := json.Marshal(cloudEvt); err == nil {
				h.hub.BroadcastToTenant(tenantID, "*", data)
			}
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(evt)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *EventHandler) HandleEventByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/events/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid event id"}`, http.StatusBadRequest)
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	if tenantID == uuid.Nil {
		tenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}

	if r.Method == http.MethodPatch || r.Method == http.MethodPost {
		_ = h.eventRepo.UpdateStatus(r.Context(), tenantID, id, domain.EventStatusAcknowledged, nil)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "acknowledged"})
		return
	}

	evt, err := h.eventRepo.GetByID(r.Context(), tenantID, id)
	if err != nil {
		http.Error(w, `{"error":"event not found"}`, http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(evt)
}
