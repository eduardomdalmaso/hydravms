package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	TenantID string   `json:"tenant_id"`
	UserID   string   `json:"user_id"`
	Role     string   `json:"role"`
	Scopes   []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

type TokenValidatorFunc func(tokenStr string) (*TokenClaims, error)

// AuthMiddleware enforces strict Token-Only access (JWT & API Keys) and injects tenant_id into context.
func AuthMiddleware(validateToken TokenValidatorFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CORS preflight, public auth login and health check endpoints
			if r.Method == http.MethodOptions || r.URL.Path == "/api/v1/auth/login" || strings.HasPrefix(r.URL.Path, "/swagger/") || r.URL.Path == "/healthz" {
				next.ServeHTTP(w, r)
				return
			}

			var rawToken string
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					rawToken = parts[1]
				}
			}
			if rawToken == "" {
				rawToken = r.URL.Query().Get("token")
			}
			if rawToken == "" && r.Header.Get("Sec-WebSocket-Protocol") != "" {
				rawToken = strings.TrimSpace(r.Header.Get("Sec-WebSocket-Protocol"))
			}

			if rawToken == "" {
				renderError(w, http.StatusUnauthorized, "missing authorization token")
				return
			}

			claims, err := validateToken(rawToken)
			if err != nil {
				renderError(w, http.StatusUnauthorized, "invalid or expired authentication token")
				return
			}

			tenantUUID, err := uuid.Parse(claims.TenantID)
			if err != nil || tenantUUID == uuid.Nil {
				renderError(w, http.StatusForbidden, "malformed tenant identity in token claims")
				return
			}

			userUUID, _ := uuid.Parse(claims.UserID)

			// Propagate tenant and user identities securely in context
			ctx := context.WithValue(r.Context(), TenantContextKey, tenantUUID)
			if userUUID != uuid.Nil {
				ctx = context.WithValue(ctx, UserContextKey, userUUID)
			}
			if claims.Role != "" {
				ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func renderError(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"type":   "https://hydravms.domain.com/errors/unauthorized",
		"title":  "Authentication Error",
		"status": status,
		"detail": detail,
	})
}
