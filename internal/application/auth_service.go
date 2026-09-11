package application

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"hydravms/internal/domain"
	"hydravms/internal/ports"
)

type JWTClaims struct {
	TenantID string   `json:"tenant_id"`
	UserID   string   `json:"user_id"`
	Role     string   `json:"role"`
	Scopes   []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo     ports.UserRepository
	auditService *AuditService
	jwtSecret    []byte
}

func NewAuthService(userRepo ports.UserRepository, auditService *AuditService) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "hydravms-super-secret-production-key-2026-auth-token-guard"
	}
	return &AuthService{
		userRepo:     userRepo,
		auditService: auditService,
		jwtSecret:    []byte(secret),
	}
}

func (s *AuthService) Authenticate(ctx context.Context, identifier, password, ipAddress, userAgent string) (string, *domain.User, error) {
	if identifier == "" || password == "" {
		return "", nil, domain.ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmailOrUsername(ctx, identifier)
	if err != nil {
		if s.auditService != nil {
			_ = s.auditService.RecordAction(
				ctx,
				"00000000-0000-0000-0000-000000000001",
				identifier,
				ipAddress,
				"USER_LOGIN_FAILED",
				"auth",
				"user_unknown",
				fmt.Sprintf("Tentativa de login falhou: usuário %s não encontrado", identifier),
				domain.AuditCategoryAudit,
				domain.AuditLevelWarning,
				map[string]interface{}{"identifier": identifier, "ip": ipAddress, "user_agent": userAgent},
			)
		}
		return "", nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", nil, domain.ErrUserInactive
	}

	// Bcrypt password hash comparison
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if s.auditService != nil {
			_ = s.auditService.RecordAction(
				ctx,
				user.TenantID.String(),
				user.Email,
				ipAddress,
				"USER_LOGIN_FAILED",
				"auth",
				user.ID.String(),
				fmt.Sprintf("Tentativa de login falhou para %s: senha incorreta", user.Email),
				domain.AuditCategoryAudit,
				domain.AuditLevelWarning,
				map[string]interface{}{"identifier": identifier, "ip": ipAddress, "user_agent": userAgent},
			)
		}
		return "", nil, domain.ErrInvalidCredentials
	}

	// Generate secure JWT token (24 hours expiry)
	tokenExpires := time.Now().Add(24 * time.Hour)
	claims := &JWTClaims{
		TenantID: user.TenantID.String(),
		UserID:   user.ID.String(),
		Role:     string(user.Role),
		Scopes:   []string{"cameras:read", "cameras:write", "events:read", "storage:read", "admin:access"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    "hydravms-controlplane",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(tokenExpires),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	if s.auditService != nil {
		_ = s.auditService.RecordAction(
			ctx,
			user.TenantID.String(),
			user.Email,
			ipAddress,
			"USER_LOGIN_SUCCESS",
			"auth",
			user.ID.String(),
			fmt.Sprintf("Usuário %s autenticado com sucesso no sistema", user.Email),
			domain.AuditCategoryAudit,
			domain.AuditLevelInfo,
			map[string]interface{}{"role": string(user.Role), "ip": ipAddress},
		)
	}

	return signedToken, user, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}
