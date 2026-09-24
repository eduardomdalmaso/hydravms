package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type PluginRepository interface {
	ListPlugins(ctx context.Context, tenantID uuid.UUID) ([]*domain.Plugin, error)
	GetPluginByID(ctx context.Context, tenantID uuid.UUID, id string) (*domain.Plugin, error)
	InstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string, version string) error
	UninstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) error
	TogglePlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) (*domain.Plugin, error)
}
