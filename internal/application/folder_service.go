package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type FolderService struct {
	repo ports.FolderRepository
}

func NewFolderService(repo ports.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule, name string, parentID *uuid.UUID, colorHex string) (*domain.Folder, error) {
	cleanName := strings.TrimSpace(strings.ToUpper(name))
	if cleanName == "" {
		return nil, fmt.Errorf("folder name cannot be empty")
	}

	if colorHex == "" {
		colorHex = "#ff5e3a"
	}

	folder := &domain.Folder{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Module:    module,
		ParentID:  parentID,
		Name:      cleanName,
		ColorHex:  colorHex,
		Icon:      "folder",
		SortOrder: 0,
	}

	if err := s.repo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *FolderService) ListFolderTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error) {
	return s.repo.ListTree(ctx, tenantID, module)
}

func (s *FolderService) DeleteFolder(ctx context.Context, tenantID, folderID uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, folderID)
}
