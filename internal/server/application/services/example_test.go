package services_test

import (
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
)

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
