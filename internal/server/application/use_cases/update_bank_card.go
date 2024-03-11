package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UpdateBankCard - Сценарий использования для обновления существующей зашифрованной банковской карты
type UpdateBankCard struct {
	// BankCardRepository - Интерфейс репозитория для сохранения банковской карты
	BankCardRepository domain.BankCardRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u UpdateBankCard) Do(
	userID, id uuid.UUID,
	number, validThru, cvv, cardHolder string,
) error {
	card, err := u.BankCardRepository.Get(userID, id)
	if err != nil {
		return err
	}
	if card == nil {
		return domain.ErrEntityNotFound
	}

	encryptedNumber, err := u.Crypto.Encrypt([]byte(number))
	if err != nil {
		return err
	}
	encryptedValidThru, err := u.Crypto.Encrypt([]byte(validThru))
	if err != nil {
		return err
	}
	encryptedCVV, err := u.Crypto.Encrypt([]byte(cvv))
	if err != nil {
		return err
	}
	encryptedCardHolder, err := u.Crypto.Encrypt([]byte(cardHolder))
	if err != nil {
		return err
	}

	card.Number = encryptedNumber
	card.ValidThru = encryptedValidThru
	card.CVV = encryptedCVV
	card.CardHolder = encryptedCardHolder

	err = u.BankCardRepository.Update(card)

	return err
}
