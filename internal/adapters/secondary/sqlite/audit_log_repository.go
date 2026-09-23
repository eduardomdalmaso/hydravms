package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type SQLiteAuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *SQLiteAuditLogRepository {
	return &SQLiteAuditLogRepository{db: db}
}

func (r *SQLiteAuditLogRepository) Save(ctx context.Context, log *domain.AuditLog) error {
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

	var tenantIDParam string
	if log.TenantID != "" && log.TenantID != "tenant_alpha" {
		tenantIDParam = log.TenantID
	} else {
		tenantIDParam = "00000000-0000-0000-0000-000000000001"
	}

	query := `
		INSERT INTO audit_logs (tenant_id, user_id, ip_address, action, entity_type, entity_id, payload_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(ctx, query,
		tenantIDParam,
		log.UserID,
		log.IPAddress,
		log.Action,
		log.EntityType,
		log.EntityID,
		string(payloadBytes),
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	id, _ := res.LastInsertId()
	log.ID = fmt.Sprintf("%d", id)
	return nil
}

func (r *SQLiteAuditLogRepository) List(ctx context.Context, tenantID string, category string, level string, limit int, offset int) ([]*domain.AuditLog, int, error) {
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

	if tenantID != "" && tenantID != "ALL" {
		whereClauses = append(whereClauses, "a.tenant_id = ?")
		args = append(args, tenantID)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) " + baseQuery + whereSQL
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Select rows
	selectQuery := fmt.Sprintf(`
		SELECT 
			a.id,
			COALESCE(a.tenant_id, '00000000-0000-0000-0000-000000000001'),
			COALESCE(t.name, 'EMPRESA ALFA // MATRIZ'),
			COALESCE(a.user_id, ''),
			COALESCE(u.name, u.email, 'SISTEMA // AUTOMATICO'),
			COALESCE(a.ip_address, '127.0.0.1'),
			a.action,
			a.entity_type,
			COALESCE(a.entity_id, ''),
			COALESCE(a.payload_json, '{}'),
			a.created_at
		%s %s
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`, baseQuery, whereSQL)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var results []*domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		var idInt int64
		var payloadStr string
		var tenantIDStr, tenantNameStr, userIDStr, actorStr, ipStr, entityIDStr string

		if err := rows.Scan(
			&idInt,
			&tenantIDStr,
			&tenantNameStr,
			&userIDStr,
			&actorStr,
			&ipStr,
			&l.Action,
			&l.EntityType,
			&entityIDStr,
			&payloadStr,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log row: %w", err)
		}

		l.ID = fmt.Sprintf("%d", idInt)
		l.TenantID = tenantIDStr
		l.TenantName = tenantNameStr
		l.UserID = userIDStr
		l.Actor = actorStr
		l.IPAddress = ipStr
		l.EntityID = entityIDStr
		l.Target = fmt.Sprintf("%s // %s", l.EntityType, l.EntityID)

		var payloadMap map[string]interface{}
		if len(payloadStr) > 0 {
			_ = json.Unmarshal([]byte(payloadStr), &payloadMap)
		}
		l.PayloadJSON = payloadMap

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

		if category != "" && category != "ALL" && string(l.Category) != category {
			continue
		}
		if level != "" && level != "ALL" && string(l.Level) != level {
			continue
		}

		results = append(results, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return results, total, nil
}

func (r *SQLiteAuditLogRepository) Count(ctx context.Context, tenantID string) (int, error) {
	query := "SELECT COUNT(*) FROM audit_logs"
	var args []interface{}
	if tenantID != "" && tenantID != "00000000-0000-0000-0000-000000000001" {
		query += " WHERE tenant_id = ?"
		args = append(args, tenantID)
	}
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var _ ports.AuditLogRepository = (*SQLiteAuditLogRepository)(nil)
