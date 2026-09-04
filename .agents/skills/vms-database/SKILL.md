---
name: vms-database
description: Relational SQL database schema, multi-tenant isolation, composite indexing, JSONB polygon/bbox querying, pgx connection pooling, and migration guidelines for HydraVMS.
---

# 🗄️ Relational SQL Database & Multi-Tenant Data Architecture (HydraVMS)

This skill documents the relational SQL schema (PostgreSQL 15+ / SQLite WAL), query patterns, and strict multi-tenant isolation rules.

---

## 1. Table Groups & Relations

- **Tenants & RBAC:** `tenants`, `users`, `user_sessions`.
- **Devices & Streaming:** `cameras`, `camera_recording_profiles`, `camera_zones`.
- **AI Analytics & Triggers:** `ai_rules`, `events`, `event_audit_actions`.
- **Storage & Recordings:** `recordings`, `storage_pools`, `storage_purge_logs`, `video_exports`.
- **Governance:** `audit_logs`.

---

## 2. Mandatory Multi-Tenant Query Rules
Every repository implementation MUST enforce the `tenant_id` WHERE clause:
```sql
SELECT * FROM events 
WHERE tenant_id = $1 
  AND triggered_at >= $2 
ORDER BY triggered_at DESC 
LIMIT 50;
```
