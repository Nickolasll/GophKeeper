// Package jose содержит имплементацию сервиса JOSE
package jose

import (
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// JOSEService - JavaScript Object Signing and Encryption Service
type JOSEService struct {
	// TokenExp - Время жизни токена в секундах
	TokenExp time.Duration
	// JWKS - Ключ подписания токена
	JWKs jwk.Key
	log  *logrus.Logger
}

// IssueToken - Выпускает JWT для userID и подписывает его с помощью jwks
func (jose JOSEService) IssueToken(userID uuid.UUID) ([]byte, error) {
	issuedAt := time.Now()
	expiration := issuedAt.Add(jose.TokenExp)
	token, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(expiration).
		Claim("UserID", userID.String()).
		Build()
	if err != nil {
		return []byte{}, err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, jose.JWKs))
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
		jwt.WithKey(jwa.HS256, jose.JWKs),
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

// GetCerts - Возвращает публичный ключ для верификации jwt на стороне клиента
func (jose JOSEService) GetCerts() (jwk.Key, error) {
	key, err := jose.JWKs.PublicKey()
	if err != nil {
		return key, err
	}

	return key, nil
}

// New - Возвращает инстанс JOSEService, проверяет jwks на валидность
func New(
	rawJWK []byte,
	jwtExpiration time.Duration,
	log *logrus.Logger,
) (*JOSEService, error) {
	key, err := jwk.FromRaw(rawJWK)
	if err != nil {
		return nil, err
	}
	joseService := JOSEService{
		JWKs:     key,
		TokenExp: jwtExpiration,
		log:      log,
	}

	return &joseService, nil
}
