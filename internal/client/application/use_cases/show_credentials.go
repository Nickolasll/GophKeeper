// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// ShowCredentials - Сценарий получения расшифрованных логина и пароля
type ShowCredentials struct {
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// CredentialsRepository - Реализация интерфейса CredentialsRepositoryInterface
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u ShowCredentials) Do(session domain.Session) ([]domain.Credentials, error) {
	result := []domain.Credentials{}
	_, err := u.CheckToken.Do(session.Token)
	if err != nil {
		return result, err
	}

	result, err = u.CredentialsRepository.GetAll(session.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
