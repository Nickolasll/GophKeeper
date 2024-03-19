package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// GetAllCredentials - Сценарий использования для получения всех расшифрованных логинов и паролей
type GetAllCredentials struct {
	// CredentialsRepository - Интерфейс репозитория для получения логинов и паролей
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Crypto - Сервис для дешифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает слайс расшифрованных логинов и паролей
func (u GetAllCredentials) Do(userID uuid.UUID) ([]*domain.Credentials, error) {
	creds, err := u.CredentialsRepository.GetAll(userID)
	if err != nil {
		return []*domain.Credentials{}, err
	}
	for i, v := range creds {
		decryptedName, err := u.Crypto.Decrypt(v.Name)
		if err != nil {
			return []*domain.Credentials{}, err
		}
		decryptedLogin, err := u.Crypto.Decrypt(v.Login)
		if err != nil {
			return []*domain.Credentials{}, err
		}
		decryptedPassword, err := u.Crypto.Decrypt(v.Password)
		if err != nil {
			return []*domain.Credentials{}, err
		}
		decryptedMeta, err := u.Crypto.Decrypt(v.Meta)
		if err != nil {
			return []*domain.Credentials{}, err
		}
		v.Name = decryptedName
		v.Login = decryptedLogin
		v.Password = decryptedPassword
		v.Meta = decryptedMeta
		creds[i] = v
	}

	return creds, nil
}
