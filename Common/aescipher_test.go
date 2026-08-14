package Common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"

	aesccm "github.com/pschlump/AesCCM"
)

func TestDecryptAESCipherToken(t *testing.T) {
	key := []byte("1234567890abcdef")
	iv := []byte("abcdefghijklmnop")
	plaintext := []byte("robot-token-value")

	encryptedToken := buildEncryptedTokenForTest(t, key, iv, plaintext)
	t.Setenv("DATABASE_SECRET_KEY", string(key))

	decrypted, err := DecryptAESCipherToken(encryptedToken)
	if err != nil {
		t.Fatalf("unexpected error while decrypting token: %v", err)
	}

	if decrypted != string(plaintext) {
		t.Fatalf("expected %q, got %q", string(plaintext), decrypted)
	}
}

func TestDecryptAESCipherTokenReturnsErrorWithInvalidKey(t *testing.T) {
	t.Setenv("DATABASE_SECRET_KEY", "")
	_, err := DecryptAESCipherToken("YWJj")
	if err == nil {
		t.Fatal("expected an error when DATABASE_SECRET_KEY is not set")
	}
}

func TestDecryptAESCipherTokenReturnsErrorWithInvalidPayload(t *testing.T) {
	t.Setenv("DATABASE_SECRET_KEY", "1234567890abcdef")
	_, err := DecryptAESCipherToken("invalid-base64")
	if err == nil {
		t.Fatal("expected an error for invalid token payload")
	}
}

func TestDecryptAESCipherTokenWithUUIDQuaySecretKey(t *testing.T) {
	quaySecret := "00112233-4455-6677-8899-aabbccddeeff"
	derivedKey, err := deriveQuaySecretKey(quaySecret)
	if err != nil {
		t.Fatalf("unexpected error while deriving key: %v", err)
	}

	iv := []byte("abcdefghijklmnop")
	plaintext := []byte("robot-token-guid")
	encryptedToken := buildEncryptedTokenForTest(t, derivedKey, iv, plaintext)

	t.Setenv("DATABASE_SECRET_KEY", quaySecret)
	decrypted, err := DecryptAESCipherToken(encryptedToken)
	if err != nil {
		t.Fatalf("unexpected error while decrypting token: %v", err)
	}

	if decrypted != string(plaintext) {
		t.Fatalf("expected %q, got %q", string(plaintext), decrypted)
	}
}

func TestDecryptAESCipherTokenWithQuayV0Format(t *testing.T) {
	quaySecret := "00112233-4455-6677-8899-aabbccddeeff"
	derivedKey, err := deriveQuaySecretKey(quaySecret)
	if err != nil {
		t.Fatalf("unexpected error while deriving key: %v", err)
	}

	plaintext := []byte("robot-token-v0")
	nonce := []byte("1234567890123")
	encryptedToken := buildVersionedV0EncryptedTokenForTest(t, derivedKey, nonce, plaintext)

	t.Setenv("DATABASE_SECRET_KEY", quaySecret)
	decrypted, err := DecryptAESCipherToken(encryptedToken)
	if err != nil {
		t.Fatalf("unexpected error while decrypting versioned token: %v", err)
	}

	if decrypted != string(plaintext) {
		t.Fatalf("expected %q, got %q", string(plaintext), decrypted)
	}
}

func buildEncryptedTokenForTest(t *testing.T, key []byte, iv []byte, plaintext []byte) string {
	t.Helper()

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("unexpected key error: %v", err)
	}

	padded := applyPKCS7PaddingForTest(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
}

func applyPKCS7PaddingForTest(data []byte, blockSize int) []byte {
	paddingLength := blockSize - (len(data) % blockSize)
	padded := make([]byte, len(data)+paddingLength)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(paddingLength)
	}

	return padded
}

func buildVersionedV0EncryptedTokenForTest(t *testing.T, key []byte, nonce []byte, plaintext []byte) string {
	t.Helper()

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("unexpected key error: %v", err)
	}

	ccmCipher, err := aesccm.NewCCM(block, 16, 13)
	if err != nil {
		t.Fatalf("unexpected ccm init error: %v", err)
	}

	ciphertextWithTag := ccmCipher.Seal(nil, nonce, plaintext, nil)
	payload := append(nonce, ciphertextWithTag...)

	return "v0$$" + base64.StdEncoding.EncodeToString(payload)
}
