---
name: vms-events
description: Technical guidelines, rules, data contracts, polygon calculation, and multi-tenant CRUD patterns for AI computer vision events in HydraVMS.
---

# 🎯 AI Event Orchestration & CRUD Skill (HydraVMS)

This skill covers the lifecycle, data structures, triggers, and multi-tenant persistence of AI events in HydraVMS.

---

## 1. Event Rule Structure

```json
{
  "rule_id": "rule_entrance_intrusion",
  "tenant_id": "tenant_alpha_01",
  "camera_id": "cam_lobby_01",
  "event_type": "intrusion",
  "target_classes": ["person", "bicycle"],
  "confidence_threshold": 0.65,
  "polygon": [
    {"x": 0.15, "y": 0.20},
    {"x": 0.85, "y": 0.20},
    {"x": 0.85, "y": 0.90},
    {"x": 0.15, "y": 0.90}
  ],
  "cooldown_seconds": 15,
  "actions": {
    "record_clip": true,
    "clip_duration_seconds": 15,
    "notify_websocket": true,
    "webhook_url": "https://api.tenant.com/alerts"
  }
}
```

---

## 2. Multi-Tenant Scoping Rule
- Every event query (`GET /api/v1/events`), rule mutation (`POST/PUT/DELETE /api/v1/rules`) MUST inject `tenant_id` from the JWT / Session claims into the database queries (`WHERE tenant_id = ?`).
