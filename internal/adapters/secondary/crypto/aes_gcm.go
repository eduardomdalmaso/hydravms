package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	EncPrefix   = "enc:gcm:"
	nonceLength = 12 // Standard GCM nonce size
)

var defaultDevSalt = []byte("hydra-vms-internal-master-key-salt-2026")

// DeriveKey derives a consistent 32-byte key using SHA-256.
func DeriveKey(rawKey []byte) []byte {
	if len(rawKey) == 32 {
		return rawKey
	}
	hash := sha256.Sum256(rawKey)
	return hash[:]
}

// GetMasterKey retrieves the encryption key from environment or fallback dev key.
func GetMasterKey() []byte {
	envKey := strings.TrimSpace(os.Getenv("HYDRA_ENCRYPTION_KEY"))
	if envKey != "" {
		if decoded, err := base64.StdEncoding.DecodeString(envKey); err == nil && len(decoded) == 32 {
			return decoded
		}
		return DeriveKey([]byte(envKey))
	}
	return DeriveKey(defaultDevSalt)
}

// Encrypt encrypts plaintext using AES-256-GCM with a fresh random 12-byte nonce.
// Returns enc:gcm:<base64(nonce + ciphertext + tag)>.
func Encrypt(plaintext string, customKey ...[]byte) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	if strings.HasPrefix(plaintext, EncPrefix) {
		return plaintext, nil
	}

	key := GetMasterKey()
	if len(customKey) > 0 && len(customKey[0]) > 0 {
		key = DeriveKey(customKey[0])
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm mode creation failed: %w", err)
	}

	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce generation failed: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return EncPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an enc:gcm: ciphertext. If plaintext is not encrypted with the prefix,
// it is returned as-is to preserve full backward compatibility with legacy credentials.
func Decrypt(encryptedText string, customKey ...[]byte) (string, error) {
	if encryptedText == "" {
		return "", nil
	}
	if !strings.HasPrefix(encryptedText, EncPrefix) {
		return encryptedText, nil // Legacy plaintext fallback
	}

	key := GetMasterKey()
	if len(customKey) > 0 && len(customKey[0]) > 0 {
		key = DeriveKey(customKey[0])
	}

	b64Payload := strings.TrimPrefix(encryptedText, EncPrefix)
	raw, err := base64.StdEncoding.DecodeString(b64Payload)
	if err != nil {
		return "", fmt.Errorf("invalid base64 ciphertext: %w", err)
	}

	if len(raw) < nonceLength {
		return "", errors.New("ciphertext too short to contain nonce")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm mode creation failed: %w", err)
	}

	nonce := raw[:nonceLength]
	actualCiphertext := raw[nonceLength:]

	plaintextBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("aes-gcm authentication decryption failed: %w", err)
	}

	return string(plaintextBytes), nil
}
