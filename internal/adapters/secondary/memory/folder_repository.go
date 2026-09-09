package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"hydravms/internal/domain"
)

// InMemoryFolderRepository provides thread-safe folder persistence with initial seed data.
type InMemoryFolderRepository struct {
	mu      sync.RWMutex
	folders map[uuid.UUID]*domain.Folder
}

func NewInMemoryFolderRepository() *InMemoryFolderRepository {
	repo := &InMemoryFolderRepository{
		folders: make(map[uuid.UUID]*domain.Folder),
	}
	repo.seedInitialFolders()
	return repo
}

func (r *InMemoryFolderRepository) seedInitialFolders() {
	defaultTenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	modules := []domain.FolderModule{
		domain.ModuleCameras,
		domain.ModuleLayouts,
		domain.ModuleMaps,
		domain.ModuleTours,
		domain.ModuleWorkflows,
		domain.ModuleUsers,
	}

	names := map[domain.FolderModule][]string{
		domain.ModuleCameras:   {"PORTARIA PRINCIPAL", "GALPAO 01", "PERIMETRO EXTERNO"},
		domain.ModuleLayouts:   {"MOSAICOS OPERACIONAIS", "GRIDS DE EMERGENCIA"},
		domain.ModuleMaps:      {"PLANTA BAIXA CENTRAL", "MAPAS GIS ESTACIONAMENTO"},
		domain.ModuleTours:     {"RONDA NOTURNA 24/7", "PATRULHA DE ENTRADA"},
		domain.ModuleWorkflows: {"ALERTAS TELEGRAM CRITICOS", "NOTIFICACOES WEBSOCKET"},
		domain.ModuleUsers:     {"OPERADORES DE TURNO", "SUPERVISORES FORENSES"},
	}

	for _, mod := range modules {
		for idx, name := range names[mod] {
			id := uuid.New()
			r.folders[id] = &domain.Folder{
				ID:        id,
				TenantID:  defaultTenantID,
				Module:    mod,
				Name:      name,
				ColorHex:  "#ff5e3a",
				Icon:      "folder",
				SortOrder: idx,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}
}

func (r *InMemoryFolderRepository) Create(ctx context.Context, folder *domain.Folder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if folder.ID == uuid.Nil {
		folder.ID = uuid.New()
	}
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	r.folders[folder.ID] = folder
	return nil
}

func (r *InMemoryFolderRepository) GetByID(ctx context.Context, tenantID, folderID uuid.UUID) (*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	folder, exists := r.folders[folderID]
	if !exists || (folder.TenantID != tenantID && folder.TenantID != uuid.MustParse("00000000-0000-0000-0000-000000000001")) {
		return nil, fmt.Errorf("folder not found")
	}

	return folder, nil
}

func (r *InMemoryFolderRepository) ListTree(ctx context.Context, tenantID uuid.UUID, module domain.FolderModule) ([]*domain.Folder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Folder
	for _, f := range r.folders {
		if (f.TenantID == tenantID || f.TenantID == uuid.MustParse("00000000-0000-0000-0000-000000000001")) && f.Module == module {
			result = append(result, f)
		}
	}

	return result, nil
}

func (r *InMemoryFolderRepository) Update(ctx context.Context, folder *domain.Folder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.folders[folder.ID]
	if !exists {
		return fmt.Errorf("folder not found")
	}

	existing.Name = folder.Name
	existing.ColorHex = folder.ColorHex
	existing.Icon = folder.Icon
	existing.ParentID = folder.ParentID
	existing.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryFolderRepository) Delete(ctx context.Context, tenantID, folderID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.folders, folderID)
	return nil
}
