package usecases

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// Login - Сценарий использования входа пользователя по логину и паролю
type Login struct {
	// UserRepository - Интерфейс репозитория пользователя
	UserRepository domain.UserRepositoryInterface
	// JOSE - Сервис выдачи и верификации JWT
	JOSE *jose.JOSEService
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u Login) Do(login, password string) ([]byte, error) {
	var token []byte

	user, err := u.UserRepository.GetByLogin(login)
	if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
		return token, err
	}
	if user == nil {
		return token, domain.ErrLoginOrPasswordIsInvalid
	}

	if !u.JOSE.VerifyPassword(user.Password, password) {
		return token, domain.ErrLoginOrPasswordIsInvalid
	}
	token, err = u.JOSE.IssueToken(user.ID)
	if err != nil {
		return token, err
	}

	return token, err
}
