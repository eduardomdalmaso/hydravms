package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type AuditLogRepository struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepository(pool *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{pool: pool}
}

func (r *AuditLogRepository) Save(ctx context.Context, log *domain.AuditLog) error {
	if log == nil {
		return fmt.Errorf("audit log cannot be nil")
	}

	payloadBytes, err := json.Marshal(log.PayloadJSON)
	if err != nil {
		payloadBytes = []byte("{}")
	}

	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	var tenantIDParam interface{}
	if log.TenantID != "" && log.TenantID != "tenant_alpha" {
		tenantIDParam = log.TenantID
	} else {
		tenantIDParam = "00000000-0000-0000-0000-000000000001"
	}

	var userIDParam interface{}
	if log.UserID != "" {
		userIDParam = log.UserID
	}

	query := `
		INSERT INTO audit_logs (tenant_id, user_id, ip_address, action, entity_type, entity_id, payload_json, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id string
	err = r.pool.QueryRow(ctx, query,
		tenantIDParam,
		userIDParam,
		log.IPAddress,
		log.Action,
		log.EntityType,
		log.EntityID,
		payloadBytes,
		log.CreatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	log.ID = id
	return nil
}

func (r *AuditLogRepository) List(ctx context.Context, tenantID string, category string, level string, limit int, offset int) ([]*domain.AuditLog, int, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	baseQuery := `
		FROM audit_logs a
		LEFT JOIN tenants t ON t.id = a.tenant_id
		LEFT JOIN users u ON u.id = a.user_id
	`
	whereClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if tenantID != "" && tenantID != "00000000-0000-0000-0000-000000000001" && tenantID != "tenant_alpha" {
		whereClauses = append(whereClauses, fmt.Sprintf("a.tenant_id = $%d", argIdx))
		args = append(args, tenantID)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) " + baseQuery + whereSQL
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Select rows
	selectQuery := fmt.Sprintf(`
		SELECT 
			a.id::text,
			COALESCE(a.tenant_id::text, '00000000-0000-0000-0000-000000000001'),
			COALESCE(t.name, 'EMPRESA ALFA // MATRIZ'),
			COALESCE(a.user_id::text, ''),
			COALESCE(u.name, u.email, 'SISTEMA // AUTOMATICO'),
			COALESCE(a.ip_address, '127.0.0.1'),
			a.action,
			a.entity_type,
			COALESCE(a.entity_id, ''),
			COALESCE(a.payload_json, '{}'::jsonb),
			a.created_at
		%s %s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, whereSQL, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var results []*domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		var payloadBytes []byte
		var tenantIDStr, tenantNameStr, userIDStr, actorStr, ipStr, entityIDStr string

		if err := rows.Scan(
			&l.ID,
			&tenantIDStr,
			&tenantNameStr,
			&userIDStr,
			&actorStr,
			&ipStr,
			&l.Action,
			&l.EntityType,
			&entityIDStr,
			&payloadBytes,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log row: %w", err)
		}

		l.TenantID = tenantIDStr
		l.TenantName = tenantNameStr
		l.UserID = userIDStr
		l.Actor = actorStr
		l.IPAddress = ipStr
		l.EntityID = entityIDStr
		l.Target = fmt.Sprintf("%s // %s", l.EntityType, l.EntityID)

		var payloadMap map[string]interface{}
		if len(payloadBytes) > 0 {
			_ = json.Unmarshal(payloadBytes, &payloadMap)
		}
		l.PayloadJSON = payloadMap

		// Derive category and level
		l.Category = domain.AuditCategoryAudit
		if strings.HasPrefix(l.Action, "SYSTEM_") || l.EntityType == "system" || l.EntityType == "watchdog" {
			l.Category = domain.AuditCategorySystem
		}

		l.Level = domain.AuditLevelInfo
		if payloadMap != nil {
			if lvl, ok := payloadMap["level"].(string); ok {
				l.Level = domain.AuditLevel(strings.ToUpper(lvl))
			}
			if msg, ok := payloadMap["message"].(string); ok && msg != "" {
				l.Details = msg
			}
		}

		if l.Details == "" {
			l.Details = fmt.Sprintf("Operação %s executada com sucesso sobre %s", l.Action, l.Target)
		}

		// Apply in-memory category / level filter if requested
		if category != "" && category != "ALL" && string(l.Category) != category {
			continue
		}
		if level != "" && level != "ALL" && string(l.Level) != level {
			continue
		}

		results = append(results, &l)
	}

	return results, total, nil
}

func (r *AuditLogRepository) Count(ctx context.Context, tenantID string) (int, error) {
	query := "SELECT COUNT(*) FROM audit_logs"
	var args []interface{}
	if tenantID != "" && tenantID != "00000000-0000-0000-0000-000000000001" {
		query += " WHERE tenant_id = $1"
		args = append(args, tenantID)
	}
	var count int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

var _ ports.AuditLogRepository = (*AuditLogRepository)(nil)
