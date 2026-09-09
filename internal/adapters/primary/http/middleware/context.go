package middleware

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	TenantContextKey contextKey = "tenant_id"
	UserContextKey   contextKey = "user_id"
	RoleContextKey   contextKey = "user_role"
)

var (
	ErrMissingTenantContext = errors.New("missing or invalid tenant context")
	ErrMissingUserContext   = errors.New("missing or invalid user context")
)

// WithTenantContext returns a copy of parent context with tenantID injected.
func WithTenantContext(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenantID)
}

// GetTenantID extracts the strongly-typed tenant UUID from context.
func GetTenantID(ctx context.Context) (uuid.UUID, error) {
	val, ok := ctx.Value(TenantContextKey).(uuid.UUID)
	if !ok || val == uuid.Nil {
		return uuid.Nil, ErrMissingTenantContext
	}
	return val, nil
}

// GetUserID extracts the user UUID from context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	val, ok := ctx.Value(UserContextKey).(uuid.UUID)
	if !ok || val == uuid.Nil {
		return uuid.Nil, ErrMissingUserContext
	}
	return val, nil
}
