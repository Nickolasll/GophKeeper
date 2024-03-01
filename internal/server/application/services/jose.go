// Package services содержит имплементацию сервисов Crypto и JOSE
package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

// JOSEService - JavaScript Object Signing and Encryption Service
type JOSEService struct {
	// TokenExp - Время жизни токена в секундах
	TokenExp time.Duration
	// JWKS - Ключ подписания токена
	JWKS jwk.Key
}

// IssueToken - Выпускает JWT для userID и подписывает его с помощью jwks
func (jose JOSEService) IssueToken(userID uuid.UUID) ([]byte, error) {
	issuedAt := time.Now()
	expiration := issuedAt.Add(jose.TokenExp)
	token, err := jwt.NewBuilder().IssuedAt(time.Now()).Expiration(expiration).Claim("UserID", userID.String()).Build()
	if err != nil {
		return []byte{}, err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, jose.JWKS))
	if err != nil {
		return []byte{}, err
	}

	return signed, nil
}

// ParseUserID - Валидирует и парсит JWT, извлекает клейм UserID и возвращает в формате UUID
func (jose JOSEService) ParseUserID(signed []byte) (uuid.UUID, error) {
	var userID uuid.UUID
	token, err := jwt.Parse(
		signed,
		jwt.WithKey(jwa.HS256, jose.JWKS),
		jwt.WithValidate(true),
	)
	if err != nil {
		return userID, err
	}
	v, _ := token.Get("UserID")
	str, _ := v.(string)
	userID, err = uuid.Parse(str)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

// Hash - Хэширует пароль
func (jose JOSEService) Hash(password string) string {
	var passwordBytes = []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPassword)
}

// VerifyPassword - Сравнивает хэш и пароль
func (jose JOSEService) VerifyPassword(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))

	return err == nil
}

func (jose JOSEService) GetCerts() (jwk.Key, error) {
	key, err := jose.JWKS.PublicKey()
	if err != nil {
		return key, err
	}

	return key, nil
}
