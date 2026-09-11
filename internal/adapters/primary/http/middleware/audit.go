package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"hydravms/internal/application"
	"hydravms/internal/domain"
)

// AuditMiddleware logs all state-mutating HTTP API calls to the forensic audit log.
func AuditMiddleware(auditService *application.AuditService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only audit mutating actions (POST, PUT, PATCH, DELETE) under /api/v1/
			method := r.Method
			path := r.URL.Path

			isMutating := method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
			isAPI := strings.HasPrefix(path, "/api/v1/")
			isLogsEndpoint := strings.Contains(path, "/system/logs") || strings.Contains(path, "/audit/logs")

			if !isMutating || !isAPI || isLogsEndpoint || auditService == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Read and restore request body for auditing
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			rec := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			// If operation was accepted/successful, write forensic audit record
			if rec.statusCode >= 200 && rec.statusCode < 400 {
				go func(m, p, ip string, b []byte) {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()

					action := determineAction(m, p)
					entityType, entityID := extractEntity(p)
					details := fmt.Sprintf("%s executado com sucesso sobre %s // %s (HTTP %d)", action, entityType, entityID, rec.statusCode)

					var payload map[string]interface{}
					if len(b) > 0 {
						_ = json.Unmarshal(b, &payload)
					}
					if payload == nil {
						payload = make(map[string]interface{})
					}
					payload["http_method"] = m
					payload["http_path"] = p
					payload["status_code"] = rec.statusCode

					_ = auditService.RecordAction(
						ctx,
						"00000000-0000-0000-0000-000000000001",
						"operador_master",
						ip,
						action,
						entityType,
						entityID,
						details,
						domain.AuditCategoryAudit,
						domain.AuditLevelInfo,
						payload,
					)
				}(method, path, r.RemoteAddr, bodyBytes)
			}
		})
	}
}

func determineAction(method, path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/"), "/")
	resource := "RESOURCE"
	if len(parts) > 0 && parts[0] != "" {
		resource = strings.ToUpper(strings.TrimSuffix(parts[0], "s"))
	}

	switch method {
	case http.MethodPost:
		if strings.Contains(path, "/probe") {
			return resource + "_PROBED"
		}
		if strings.Contains(path, "/drain") {
			return resource + "_DRAINED"
		}
		return resource + "_CREATED"
	case http.MethodPut, http.MethodPatch:
		return resource + "_UPDATED"
	case http.MethodDelete:
		return resource + "_DELETED"
	default:
		return resource + "_MUTATED"
	}
}

func extractEntity(path string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/"), "/")
	if len(parts) == 0 {
		return "system", "general"
	}
	entityType := parts[0]
	entityID := ""
	if len(parts) > 1 {
		entityID = parts[1]
	}
	return entityType, entityID
}
