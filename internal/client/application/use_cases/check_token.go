package usecases

import (
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
type CheckToken struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// JWKRepository - Реализация интерфейса JWKRepositoryInterface
	JWKRepository domain.JWKRepositoryInterface
	// Key - Ключ для проверки JWT
	Key jwk.Key
	// Log - логгер
	Log *logrus.Logger
}

func (u CheckToken) getUserID(token string) (string, error) {
	tok, err := jwt.Parse(
		[]byte(token),
		jwt.WithKey(jwa.HS256, u.Key),
		jwt.WithValidate(true),
	)
	if err != nil {
		return "", err
	}
	v, _ := tok.Get("UserID")
	str, _ := v.(string)

	return str, nil
}

func (u *CheckToken) setupKey() error {
	certs, err := u.Client.GetCerts()
	if err != nil {
		return err
	}
	u.Key, err = jwk.ParseKey(certs)
	if err != nil {
		return err
	}

	if err := u.JWKRepository.Save(u.Key); err != nil {
		return err
	}

	return nil
}

func (u CheckToken) parseID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// Do - Вызов логики сценария использования
func (u CheckToken) Do(token string) (uuid.UUID, error) {
	var uid uuid.UUID
	if u.Key == nil {
		err := u.setupKey()
		if err != nil {
			return uid, err
		}
	}

	id, err := u.getUserID(token)

	if err != nil {
		err = u.setupKey()
		if err != nil {
			return uid, err
		}
		id, err = u.getUserID(token)
		if err != nil {
			return uid, domain.ErrInvalidToken
		}

		uid, err = u.parseID(id)
		if err != nil {
			return uid, err
		}

		return uid, nil
	}

	uid, err = u.parseID(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

func (u CheckToken) getSessionFromToken(token string) (domain.Session, error) {
	var session domain.Session
	userID, err := u.Do(token)
	if err != nil {
		return session, err
	}

	session = domain.Session{
		UserID: userID,
		Token:  token,
	}

	return session, nil
}
