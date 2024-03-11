package usecases

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// Registration - Сценарий использования регистрации пользователя по логину и паролю
type Registration struct {
	// UserRepository - Интерфейс репозитория пользователя
	UserRepository domain.UserRepositoryInterface
	// JOSE - Сервис выдачи и верификации JWT
	JOSE *jose.JOSEService
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u Registration) Do(login, password string) ([]byte, error) {
	var token []byte

	user, err := u.UserRepository.GetByLogin(login)
	if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
		return token, err
	}
	if user != nil {
		return token, domain.ErrLoginAlreadyInUse
	}

	hashedPassword := u.JOSE.Hash(password)
	newUser := domain.User{
		ID:       uuid.New(),
		Login:    login,
		Password: hashedPassword,
	}
	err = u.UserRepository.Create(newUser)
	if err != nil {
		return token, err
	}

	token, err = u.JOSE.IssueToken(newUser.ID)
	if err != nil {
		return token, err
	}

	return token, err
}
