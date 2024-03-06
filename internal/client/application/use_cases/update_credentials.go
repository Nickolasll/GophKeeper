// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// UpdateCredentials - Сценарий обновления существующего логина и пароля
type UpdateCredentials struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// CredentialsRepository - Реализация интерфейса CredentialsRepositoryInterface
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u UpdateCredentials) Do(session domain.Session, credID uuid.UUID, name, login, password string) error {
	cred, err := u.CredentialsRepository.Get(session.UserID, credID)
	if err != nil {
		return err
	}

	if name != "" {
		cred.Name = name
	}
	if login != "" {
		cred.Login = login
	}
	if password != "" {
		cred.Password = password
	}

	if err := u.Client.UpdateCredentials(session, cred); err != nil {
		return err
	}

	if err := u.CredentialsRepository.Update(session.UserID, cred); err != nil {
		return err
	}

	return nil
}
