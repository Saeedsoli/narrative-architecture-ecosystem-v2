// apps/backend/pkg/crypto/encryption.go

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt داده‌ها را با استفاده از کلید ارائه شده و الگوریتم AES-GCM رمزنگاری می‌کند.
func Encrypt(data, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Nonce در ابتدای ciphertext قرار می‌گیرد.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt داده‌های رمزنگاری شده با AES-GCM را رمزگشایی می‌کند.
func Decrypt(data, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("decryption key must be 32 bytes for AES-256")
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}