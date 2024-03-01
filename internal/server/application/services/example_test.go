package services_test

import (
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
)

func ExampleCryptoService() {
	cryptoService := services.CryptoService{
		SecretKey: []byte("1234567812345678"),
	}

	data := "Example string for encryption"

	encrypted, err := cryptoService.Encrypt(data)
	if err != nil {
		// Возникла проблема в процессе шифрования данных
	}

	decrypted, err := cryptoService.Decrypt(encrypted)

	if err != nil {
		// Возникла проблема в процессе дешифрования данных
	}

	if decrypted == data {
		// Данные расшифрованы успешно
	}
}

func ExampleJOSEService() {
	rawJWK := []byte("My secret keys")
	key, err := jwk.FromRaw(rawJWK)
	if err != nil {
		// В процессе формирования ключей возникла ошибка
	}
	joseService := services.JOSEService{
		TokenExp: 60,
		JWKS:     key,
	}

	userID := uuid.New()
	token, err := joseService.IssueToken(userID)
	if err != nil {
		// В процессе генерации JWT возникла ошибка
	}

	id, err := joseService.ParseUserID(token)
	if err != nil {
		// В процессе верификации токена и получения идентификатора возникла ошибка
	}

	if id == userID {
		// Успешно расшифровали токен
	}

	password := "my password"
	hash := joseService.Hash(password)
	if joseService.VerifyPassword(hash, "my password") == true {
		// Правильный пароль
	}
}
