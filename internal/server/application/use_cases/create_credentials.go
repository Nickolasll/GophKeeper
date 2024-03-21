// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// CreateCredentials - Сценарий использования для создания зашифрованного логина и пароля
type CreateCredentials struct {
	// CredentialsRepository - Интерфейс репозитория для сохранения логина и пароля
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает идентификатор ресурса
func (u *CreateCredentials) Do(
	userID uuid.UUID,
	name, login, password, meta string,
) (uuid.UUID, error) {
	credID := uuid.New()
	encryptedName, err := u.Crypto.Encrypt([]byte(name))
	if err != nil {
		return credID, err
	}
	encryptedLogin, err := u.Crypto.Encrypt([]byte(login))
	if err != nil {
		return credID, err
	}
	encryptedPassword, err := u.Crypto.Encrypt([]byte(password))
	if err != nil {
		return credID, err
	}
	encryptedMeta, err := u.Crypto.Encrypt([]byte(meta))
	if err != nil {
		return credID, err
	}
	cred := domain.Credentials{
		ID:       credID,
		UserID:   userID,
		Name:     encryptedName,
		Login:    encryptedLogin,
		Password: encryptedPassword,
		Meta:     encryptedMeta,
	}
	err = u.CredentialsRepository.Create(&cred)

	return credID, err
}
