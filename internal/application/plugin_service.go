package application

import (
	"context"
	"fmt"
	"hydravms/internal/domain"
	"hydravms/internal/ports"

	"github.com/google/uuid"
)

type PluginService struct {
	repo ports.PluginRepository
}

func NewPluginService(repo ports.PluginRepository) *PluginService {
	return &PluginService{repo: repo}
}

func (s *PluginService) ListPlugins(ctx context.Context, tenantID uuid.UUID) ([]*domain.Plugin, error) {
	return s.repo.ListPlugins(ctx, tenantID)
}

func (s *PluginService) GetPluginByID(ctx context.Context, tenantID uuid.UUID, id string) (*domain.Plugin, error) {
	return s.repo.GetPluginByID(ctx, tenantID, id)
}

func (s *PluginService) InstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) error {
	p, err := s.repo.GetPluginByID(ctx, tenantID, pluginID)
	if err != nil {
		return fmt.Errorf("plugin not found: %w", err)
	}
	return s.repo.InstallPlugin(ctx, tenantID, pluginID, p.Version)
}

func (s *PluginService) UninstallPlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) error {
	return s.repo.UninstallPlugin(ctx, tenantID, pluginID)
}

func (s *PluginService) TogglePlugin(ctx context.Context, tenantID uuid.UUID, pluginID string) (*domain.Plugin, error) {
	return s.repo.TogglePlugin(ctx, tenantID, pluginID)
}
