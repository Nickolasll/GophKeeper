package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// GetAllBankCards - Сценарий использования для получения всех расшифрованных банковских карт
type GetAllBankCards struct {
	// BankCardRepository - Интерфейс репозитория для получения банковских карт
	BankCardRepository domain.BankCardRepositoryInterface
	// Crypto - Сервис для дешифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает слайс расшифрованных логинов и паролей
func (u GetAllBankCards) Do(userID uuid.UUID) ([]*domain.BankCard, error) {
	cards, err := u.BankCardRepository.GetAll(userID)
	if err != nil {
		return []*domain.BankCard{}, err
	}
	for i, v := range cards {
		decryptedNumber, err := u.Crypto.Decrypt(v.Number)
		if err != nil {
			return []*domain.BankCard{}, err
		}
		decryptedValidThru, err := u.Crypto.Decrypt(v.ValidThru)
		if err != nil {
			return []*domain.BankCard{}, err
		}
		decryptedCVV, err := u.Crypto.Decrypt(v.CVV)
		if err != nil {
			return []*domain.BankCard{}, err
		}
		decryptedCardHolder, err := u.Crypto.Decrypt(v.CardHolder)
		if err != nil {
			return []*domain.BankCard{}, err
		}
		decryptedMeta, err := u.Crypto.Decrypt(v.Meta)
		if err != nil {
			return []*domain.BankCard{}, err
		}
		v.Number = decryptedNumber
		v.ValidThru = decryptedValidThru
		v.CVV = decryptedCVV
		v.CardHolder = decryptedCardHolder
		v.Meta = decryptedMeta
		cards[i] = v
	}

	return cards, nil
}
