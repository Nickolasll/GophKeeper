// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
type CheckToken struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client        domain.GophKeeperClientInterface
	JWKRepository domain.JWKRepositoryInterface
	Key           jwk.Key
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

// Execute - Вызов логики сценария использования
func (u CheckToken) Execute(token string) (string, error) {
	if u.Key == nil {
		err := u.setupKey()
		if err != nil {
			return "", err
		}
	}

	userID, err := u.getUserID(token)

	if err != nil {
		err := u.setupKey()
		if err != nil {
			return "", err
		}
		userID, err = u.getUserID(token)
		if err != nil {
			return "", domain.ErrInvalidToken
		}

		return userID, nil
	}

	return userID, nil
}

func (u CheckToken) getSessionFromToken(token string) (domain.Session, error) {
	var session domain.Session
	userID, err := u.Execute(token)
	if err != nil {
		return session, err
	}

	session = domain.Session{
		UserID: userID,
		Token:  token,
	}

	return session, nil
}
