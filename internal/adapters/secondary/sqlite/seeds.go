package sqlite

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func seedDefaultData(db *sql.DB) error {
	defaultTenantID := "00000000-0000-0000-0000-000000000001"
	adminID := "00000000-0000-0000-0000-000000000002"

	_, err := db.Exec(`
		INSERT OR IGNORE INTO tenants (id, slug, name, plan, is_active, created_at, updated_at)
		VALUES (?, 'master', 'Empresa Alfa Matriz', 'enterprise', 1, datetime('now'), datetime('now'))
	`, defaultTenantID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES (?, ?, 'Admin', 'admin@hydravms.io', ?, 'admin', 1, datetime('now'), datetime('now'))
	`, adminID, defaultTenantID, string(hashedPassword))
	if err != nil {
		return err
	}

	_, _ = db.Exec(`
		INSERT OR IGNORE INTO plugins (
			id, name, version, author, category, runtime, entrypoint, min_vms_version,
			permissions, config_schema, ui_schema, is_official, is_deprecated,
			package_url, package_checksum, hardware_req, description, created_at, updated_at
		) VALUES (
			'object_detection_sota',
			'DETECCAO DE OBJETOS',
			'1.0.0',
			'Hydra Vision AI',
			'analytics',
			'binary_elf',
			'main',
			'1.0.0',
			'["video:shm_read", "events:write", "cuda:ipc"]',
			'{}',
			'{}',
			1,
			0,
			'https://huggingface.co/hydravision/yolo26m-object',
			'sha256:auto',
			'CUDA 13.3 // RTX 5090',
			'Deteccao e rastreamento SOTA de pessoas, celular, animais e veiculos com Auto-SAHI, Motion-Gating e Poligono/Linha de contagem.',
			datetime('now'),
			datetime('now')
		);

		INSERT OR IGNORE INTO plugins (
			id, name, version, author, category, runtime, entrypoint, min_vms_version,
			permissions, config_schema, ui_schema, is_official, is_deprecated,
			package_url, package_checksum, hardware_req, description, created_at, updated_at
		) VALUES (
			'facial_recognition_sota',
			'RECONHECIMENTO FACIAL',
			'1.0.0',
			'Hydra Vision AI',
			'analytics',
			'python3',
			'face_analytic/worker_cli.py',
			'1.0.0',
			'["video:shm_read", "events:write", "cuda:ipc"]',
			'{"min_height": 40, "max_height": 0, "min_faces": 3, "blur_threshold": 60.0, "unknown_score": 0.60, "recognition_threshold": 0.75}',
			'{}',
			1,
			0,
			'https://huggingface.co/hydravision/facial-recognition',
			'sha256:auto',
			'CUDA 13.3 // RTX 5090',
			'Analitico Facial Forense com Person-Gating, selecao de Best Frame, Blur Laplaciano adaptativo e ReID Cross-Camera.',
			datetime('now'),
			datetime('now')
		);
	`)

	return nil
}
