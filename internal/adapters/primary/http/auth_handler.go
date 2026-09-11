package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"hydravms/internal/application"
	"hydravms/internal/domain"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  UserResponseDTO `json:"user"`
}

type UserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed","message":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "bad_request",
			"message": "JSON de login inválido",
		})
		return
	}

	identifier := strings.TrimSpace(req.Username)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Email)
	}

	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	userAgent := r.Header.Get("User-Agent")

	token, user, err := h.authService.Authenticate(r.Context(), identifier, req.Password, clientIP, userAgent)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := "Usuário ou senha incorretos."
		if errors.Is(err, domain.ErrUserInactive) {
			msg = "Conta de usuário inativa ou bloqueada."
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "invalid_credentials",
			"message": msg,
		})
		return
	}

	usernameDisplay := user.Email
	if user.Name != "" {
		usernameDisplay = user.Name
	}

	resp := LoginResponse{
		Token: token,
		User: UserResponseDTO{
			ID:       user.ID.String(),
			Username: usernameDisplay,
			Name:     user.Name,
			Email:    user.Email,
			Role:     string(user.Role),
			TenantID: user.TenantID.String(),
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
