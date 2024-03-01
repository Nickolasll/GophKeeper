// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"errors"

	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// Login - Сценарий использования входа пользователя по логину и паролю
type Login struct {
	// UserRepository - Интерфейс репозитория пользователя
	UserRepository domain.UserRepositoryInterface
	// JOSE - Сервис выдачи и верификации JWT
	JOSE services.JOSEService
}

// Execute - Вызов исполнения сценария использования
func (u Login) Execute(login, password string) ([]byte, error) {
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
