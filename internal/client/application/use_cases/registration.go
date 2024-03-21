package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// Registration - Сценарий регистрации по логину и паролю
type Registration struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// SessionRepository - Реализация интерфейса SessionRepositoryInterface
	SessionRepository domain.SessionRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u Registration) Do(login, password string) (domain.Session, error) {
	var session domain.Session
	token, err := u.Client.Register(login, password)
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
