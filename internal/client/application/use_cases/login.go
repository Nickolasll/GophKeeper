// Package usecases содержит имплементацию бизнес логики приложения
package usecases //nolint: dupl

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
func (u Login) Execute(login, password string) (domain.Session, error) {
	var session domain.Session
	token, err := u.Client.Login(login, password)
	if err != nil {
		return session, err
	}

	session, err = u.CheckToken.getSessionFromToken(token)
	if err != nil {
		return session, err
	}

	if err := u.SessionRepository.Save(session); err != nil {
		return session, err
	}

	return session, nil
}
