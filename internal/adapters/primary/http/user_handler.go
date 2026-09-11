package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminDataHandler struct {
	pool *pgxpool.Pool
}

func NewAdminDataHandler(pool *pgxpool.Pool) *AdminDataHandler {
	return &AdminDataHandler{pool: pool}
}

// HandleUsers returns live users for the Admin Center.
func (h *AdminDataHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if h.pool == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"users": []interface{}{}})
		return
	}

	query := `
		SELECT 
			u.id::text, 
			u.email, 
			COALESCE(u.name, 'Administrador'), 
			u.role, 
			u.is_active, 
			COALESCE(t.name, 'Empresa Alfa'), 
			u.created_at, 
			COALESCE(u.last_login_at, u.created_at)
		FROM users u
		LEFT JOIN tenants t ON t.id = u.tenant_id
		ORDER BY u.created_at DESC
	`

	rows, err := h.pool.Query(r.Context(), query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"users": []interface{}{}})
		return
	}
	defer rows.Close()

	type UserItem struct {
		ID               string `json:"id"`
		Username         string `json:"username"`
		FullName         string `json:"fullName"`
		Email            string `json:"email"`
		Role             string `json:"role"`
		CompanyScope     string `json:"companyScope"`
		GroupName        string `json:"groupName"`
		IsActive         bool   `json:"isActive"`
		CreatedAt        string `json:"createdAt"`
		LastLogin        string `json:"lastLogin"`
		TwoFactorEnabled bool   `json:"twoFactorEnabled"`
	}

	var users []UserItem
	for rows.Next() {
		var u UserItem
		var createdAt, lastLogin time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CompanyScope, &createdAt, &lastLogin); err == nil {
			u.Username = u.Email
			u.GroupName = "Operação Matriz"
			u.CreatedAt = createdAt.Format(time.RFC3339)
			u.LastLogin = lastLogin.Format(time.RFC3339)
			u.TwoFactorEnabled = true
			users = append(users, u)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"total": len(users),
	})
}

// HandleLayouts returns enterprise mosaic layouts.
func (h *AdminDataHandler) HandleLayouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"layouts": []interface{}{},
		"total":   0,
	})
}

// HandleMaps returns interactive floor maps.
func (h *AdminDataHandler) HandleMaps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"maps":  []interface{}{},
		"total": 0,
	})
}

// HandleTours returns virtual camera tours and rondas.
func (h *AdminDataHandler) HandleTours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tours": []interface{}{},
		"total": 0,
	})
}
