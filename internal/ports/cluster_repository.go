package ports

import (
	"context"
	"hydravms/internal/domain"

	"github.com/google/uuid"
)

type ClusterRepository interface {
	RegisterNode(ctx context.Context, node *domain.ClusterNode) error
	GetNodeByID(ctx context.Context, nodeID uuid.UUID) (*domain.ClusterNode, error)
	ListNodes(ctx context.Context, tenantID *uuid.UUID, role *domain.NodeRole) ([]*domain.ClusterNode, error)
	UpdateHeartbeat(ctx context.Context, nodeID uuid.UUID, cpu, ram, gpu float64, activeStreams int) error
	UpdateNode(ctx context.Context, node *domain.ClusterNode) error
	DeleteNode(ctx context.Context, nodeID uuid.UUID) error
}
