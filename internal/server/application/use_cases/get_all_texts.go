package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// GetAllTexts - Сценарий использования для получения всех расшифрованных текстовых данных
type GetAllTexts struct {
	// TextRepository - Интерфейс репозитория для получения текстовых данных
	TextRepository domain.TextRepositoryInterface
	// Crypto - Сервис для дешифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает слайс расшифрованных текстовых данных
func (u GetAllTexts) Do(userID uuid.UUID) ([]domain.Text, error) {
	texts, err := u.TextRepository.GetAll(userID)
	if err != nil {
		return []domain.Text{}, err
	}
	for i, v := range texts {
		decryptedContent, err := u.Crypto.Decrypt(v.Content)
		if err != nil {
			return []domain.Text{}, err
		}
		v.Content = decryptedContent
		texts[i] = v
	}

	return texts, nil
}
