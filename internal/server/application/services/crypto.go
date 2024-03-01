// Package services содержит имплементацию сервисов Crypto и JOSE
package services

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

// Encrypt - Шифрует текстовые данные, возвращает результат в формате строки
func (c CryptoService) Encrypt(value string) ([]byte, error) {
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

	result := aesgcm.Seal(nonce, nonce, []byte(value), nil)

	return result, nil
}

// Decrypt - Расшифровывает текстовые данные, возвращает результат в формате строки
func (c CryptoService) Decrypt(value []byte) (string, error) {
	aesblock, err := aes.NewCipher(c.SecretKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := value[:nonceSize], value[nonceSize:]

	result, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
