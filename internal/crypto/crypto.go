// Package crypto содержит имплементацию сервиса шифрования
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// CryptoService - Сервис для шифрования и дешифрования данных
type CryptoService struct {
	// SecretKey - Приватный ключ шифрования
	SecretKey []byte
}

func (c CryptoService) generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Encrypt - Защифровывает бинарные данные
func (c CryptoService) Encrypt(value []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(c.SecretKey)
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return []byte{}, err
	}

	nonce, err := c.generateRandom(aesgcm.NonceSize())
	if err != nil {
		return []byte{}, err
	}

	result := aesgcm.Seal(nonce, nonce, value, nil)

	return result, nil
}

// Decrypt - Расшифровывает бинарные
func (c CryptoService) Decrypt(value []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(c.SecretKey)
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return []byte{}, err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := value[:nonceSize], value[nonceSize:]

	result, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

// New инициализирует сервис, проверяет введенный ключ
func New(secretKey []byte) (*CryptoService, error) {
	aesblock, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	_, err = cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	crypto := CryptoService{
		SecretKey: secretKey,
	}

	return &crypto, nil
}
