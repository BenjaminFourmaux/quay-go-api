package Common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"quay-go-api/Common/Errors"
	"strings"

	"github.com/google/uuid"
	aesccm "github.com/pschlump/AesCCM"
)

// DecryptAESCipherToken decrypts a base64 token encoded as IV + ciphertext
// using AES-CBC with PKCS7 padding, compatible with the Python AESCipher class.
func DecryptAESCipherToken(encryptedToken string) (string, error) {
	// Retrieve the AES key from the environment variable and decode it
	encryptionKey := os.Getenv("DATABASE_SECRET_KEY")
	if encryptionKey == "" {
		return "", Errors.QuayEncryptionKeyNotSet()
	}

	if strings.Contains(encryptedToken, "$$") {
		return decryptQuayVersionedToken(encryptedToken, encryptionKey)
	}

	return decryptLegacyAesCbcToken(encryptedToken, encryptionKey)
}

func decryptQuayVersionedToken(encryptedToken string, encryptionKey string) (string, error) {
	versionPrefix, payload, hasVersion := strings.Cut(encryptedToken, "$$")
	if !hasVersion || payload == "" {
		return "", fmt.Errorf("invalid encrypted token format")
	}

	if versionPrefix != "v0" {
		return "", fmt.Errorf("unsupported encrypted token version: %s", versionPrefix)
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("invalid encrypted token payload encoding: %w", err)
	}
	if len(decodedPayload) <= 13 {
		return "", fmt.Errorf("invalid encrypted token payload length")
	}

	key, err := deriveQuaySecretKey(encryptionKey)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("invalid AES key: %w", err)
	}
	ccmCipher, err := aesccm.NewCCM(block, 16, 13)
	if err != nil {
		return "", fmt.Errorf("unable to initialize AES-CCM: %w", err)
	}

	nonce := decodedPayload[:13]
	ciphertextWithTag := decodedPayload[13:]
	plaintext, err := ccmCipher.Open(nil, nonce, ciphertextWithTag, nil)
	if err != nil {
		return "", fmt.Errorf("unable to decrypt versioned token: %w", err)
	}

	return string(plaintext), nil
}

func decryptLegacyAesCbcToken(encryptedToken string, encryptionKey string) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", fmt.Errorf("invalid encrypted token encoding: %w", err)
	}

	if len(encryptedData) < aes.BlockSize {
		return "", fmt.Errorf("invalid encrypted token length: payload is shorter than IV")
	}

	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid encrypted token length: ciphertext must be a non-empty AES block multiple")
	}

	keys, err := buildCandidateAESKeys(encryptionKey)
	if err != nil {
		return "", err
	}

	var lastErr error
	for _, key := range keys {
		block, keyErr := aes.NewCipher(key)
		if keyErr != nil {
			lastErr = keyErr
			continue
		}

		decryptedPaddedData := make([]byte, len(ciphertext))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(decryptedPaddedData, ciphertext)

		decryptedData, unpadErr := pkcs7Unpad(decryptedPaddedData, aes.BlockSize)
		if unpadErr == nil {
			return string(decryptedData), nil
		}

		lastErr = unpadErr
	}

	return "", fmt.Errorf("unable to decrypt token with provided key: %w", lastErr)
}

func buildCandidateAESKeys(rawKey string) ([][]byte, error) {
	candidates := make([][]byte, 0, 4)
	seen := make(map[string]bool)

	appendIfValid := func(key []byte) {
		if !isValidAESKeyLength(len(key)) {
			return
		}

		fingerprint := string(key)
		if seen[fingerprint] {
			return
		}

		seen[fingerprint] = true
		candidates = append(candidates, key)
	}

	if isValidAESKeyLength(len(rawKey)) {
		appendIfValid([]byte(rawKey))
	}

	if decoded, err := base64.StdEncoding.DecodeString(rawKey); err == nil {
		appendIfValid(decoded)
	}

	if decoded, err := hex.DecodeString(rawKey); err == nil {
		appendIfValid(decoded)
	}

	quayDerived, err := deriveQuaySecretKey(rawKey)
	if err == nil {
		appendIfValid(quayDerived)
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("invalid AES key: no valid key candidate could be derived from provided value")
	}

	return candidates, nil
}

func isValidAESKeyLength(keyLength int) bool {
	return keyLength == 16 || keyLength == 24 || keyLength == 32
}

func deriveQuaySecretKey(configSecretKey string) ([]byte, error) {
	secretBytes := []byte{}

	bigInt := new(big.Int)
	if _, ok := bigInt.SetString(configSecretKey, 10); ok {
		hexValue := fmt.Sprintf("%02x", bigInt)
		decoded, err := hex.DecodeString(hexValue)
		if err == nil {
			secretBytes = decoded
		}
	}

	if len(secretBytes) == 0 {
		if parsedUUID, err := uuid.Parse(configSecretKey); err == nil {
			secretBytes = parsedUUID[:]
		}
	}

	if len(secretBytes) == 0 {
		secretBytes = []byte(configSecretKey)
	}

	if len(secretBytes) == 0 {
		return nil, fmt.Errorf("empty key")
	}

	derived := make([]byte, 32)
	for i := 0; i < len(derived); i++ {
		derived[i] = secretBytes[i%len(secretBytes)]
	}

	return derived, nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid PKCS7 padded data length")
	}

	paddingLength := int(data[len(data)-1])
	if paddingLength == 0 || paddingLength > blockSize || paddingLength > len(data) {
		return nil, fmt.Errorf("invalid PKCS7 padding size")
	}

	for i := 0; i < paddingLength; i++ {
		if int(data[len(data)-1-i]) != paddingLength {
			return nil, fmt.Errorf("invalid PKCS7 padding content")
		}
	}

	return data[:len(data)-paddingLength], nil
}
