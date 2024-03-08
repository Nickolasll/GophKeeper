package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// GetAllBinaries - Сценарий использования для получения всех расшифрованных бинарных данных
type GetAllBinaries struct {
	// BinaryRepository - Интерфейс репозитория для получения бинарных данных
	BinaryRepository domain.BinaryRepositoryInterface
	// Crypto - Сервис для дешифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает слайс расшифрованных бинарных данных
func (u GetAllBinaries) Do(userID uuid.UUID) ([]domain.Binary, error) {
	bins, err := u.BinaryRepository.GetAll(userID)
	if err != nil {
		return []domain.Binary{}, err
	}
	for i, v := range bins {
		decryptedContent, err := u.Crypto.Decrypt(v.Content)
		if err != nil {
			return []domain.Binary{}, err
		}
		v.Content = decryptedContent
		bins[i] = v
	}

	return bins, nil
}
