package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CreateCredentials - Сценарий создания нового логина и пароля
type CreateCredentials struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// CredentialsRepository - Реализация интерфейса CredentialsRepositoryInterface
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u CreateCredentials) Do(
	session domain.Session,
	name, login, password, meta string,
) error {
	credID, err := u.Client.CreateCredentials(session, name, login, password, meta)
	if err != nil {
		return err
	}

	cred := domain.Credentials{
		ID:       credID,
		Name:     name,
		Login:    login,
		Password: password,
		Meta:     meta,
	}

	if err := u.CredentialsRepository.Create(session.UserID, &cred); err != nil {
		return err
	}

	return nil
}
