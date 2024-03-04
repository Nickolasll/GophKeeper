package crypto_test

import (
	"github.com/Nickolasll/goph-keeper/internal/crypto"
)

func ExampleCryptoService() {
	cryptoService := crypto.CryptoService{
		SecretKey: []byte("1234567812345678"),
	}

	data := "Example string for encryption"

	encrypted, err := cryptoService.Encrypt([]byte(data))
	if err != nil {
		// Возникла проблема в процессе шифрования данных
	}

	decrypted, err := cryptoService.Decrypt(encrypted)

	if err != nil {
		// Возникла проблема в процессе дешифрования данных
	}

	if string(decrypted) == data {
		// Данные расшифрованы успешно
	}
}
