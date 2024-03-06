// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// CreateBankCard - Сценарий использования для создания зашифрованной банковской карты
type CreateBankCard struct {
	// BankCardRepository - Интерфейс репозитория для сохранения банковской карты
	BankCardRepository domain.BankCardRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает идентификатор ресурса
func (u *CreateBankCard) Do(userID uuid.UUID, number, validThru, cvv, cardHolder string) (uuid.UUID, error) {
	cardID := uuid.New()
	encryptedNumber, err := u.Crypto.Encrypt([]byte(number))
	if err != nil {
		return cardID, err
	}
	encryptedValidThru, err := u.Crypto.Encrypt([]byte(validThru))
	if err != nil {
		return cardID, err
	}
	encryptedCVV, err := u.Crypto.Encrypt([]byte(cvv))
	if err != nil {
		return cardID, err
	}
	encryptedCardHolder, err := u.Crypto.Encrypt([]byte(cardHolder))
	if err != nil {
		return cardID, err
	}
	card := domain.BankCard{
		ID:         cardID,
		UserID:     userID,
		Number:     encryptedNumber,
		ValidThru:  encryptedValidThru,
		CVV:        encryptedCVV,
		CardHolder: encryptedCardHolder,
	}
	err = u.BankCardRepository.Create(&card)

	return cardID, err
}
