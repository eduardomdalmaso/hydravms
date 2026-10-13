package crypto

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plain := "SecretPass123!@#"
	enc, err := Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if !strings.HasPrefix(enc, EncPrefix) {
		t.Fatalf("Expected prefix %s, got %s", EncPrefix, enc)
	}

	// Should not be equal to plain
	if enc == plain {
		t.Fatalf("Encrypted text should differ from plain")
	}

	dec, err := Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if dec != plain {
		t.Fatalf("Expected %s, got %s", plain, dec)
	}
}

func TestEncrypt_IdempotentIfAlreadyEncrypted(t *testing.T) {
	plain := "mypassword"
	enc1, _ := Encrypt(plain)
	enc2, _ := Encrypt(enc1)

	if enc1 != enc2 {
		t.Fatalf("Encrypting already encrypted text should be idempotent, got %s != %s", enc1, enc2)
	}
}

func TestDecrypt_LegacyPlaintextFallback(t *testing.T) {
	legacy := "admin12345"
	dec, err := Decrypt(legacy)
	if err != nil {
		t.Fatalf("Decrypt of legacy plaintext should not error: %v", err)
	}
	if dec != legacy {
		t.Fatalf("Expected legacy string %s, got %s", legacy, dec)
	}
}

func TestDecrypt_EmptyString(t *testing.T) {
	dec, err := Decrypt("")
	if err != nil || dec != "" {
		t.Fatalf("Expected empty string, got %s, err: %v", dec, err)
	}
}

func TestDecrypt_CorruptedCiphertext(t *testing.T) {
	corrupted := EncPrefix + "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="
	_, err := Decrypt(corrupted)
	if err == nil {
		t.Fatalf("Expected decryption error for corrupted ciphertext")
	}
}
