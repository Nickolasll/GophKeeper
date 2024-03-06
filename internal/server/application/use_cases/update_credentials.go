// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UpdateCredentials - Сценарий использования для обновления существующих зашифрованных логина и пароля
type UpdateCredentials struct {
	// CredentialstRepository - Интерфейс репозитория для сохранения логина и пароля
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u UpdateCredentials) Do(userID, id uuid.UUID, name, login, password string) error {
	cred, err := u.CredentialsRepository.Get(userID, id)
	if err != nil {
		return err
	}
	if cred == nil {
		return domain.ErrEntityNotFound
	}

	encryptedName, err := u.Crypto.Encrypt([]byte(name))
	if err != nil {
		return err
	}
	encryptedLogin, err := u.Crypto.Encrypt([]byte(login))
	if err != nil {
		return err
	}
	encryptedPassword, err := u.Crypto.Encrypt([]byte(password))
	if err != nil {
		return err
	}
	cred.Name = encryptedName
	cred.Login = encryptedLogin
	cred.Password = encryptedPassword
	err = u.CredentialsRepository.Update(cred)

	return err
}
