// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// Login - Сценарий входа по логину и паролю
type Login struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// SessionRepository - Реализация интерфейса SessionRepository
	SessionRepository domain.SessionRepositoryInterface
}

// Execute - Вызов логики сценария использования
func (u Login) Execute(login, password string) error {
	token, err := u.Client.Login(login, password)
	if err != nil {
		return err
	}

	session, err := u.CheckToken.getSessionFromToken(token)
	if err != nil {
		return err
	}

	if err := u.SessionRepository.Save(session); err != nil {
		return err
	}

	return nil
}
