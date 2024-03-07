// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// UpdateBankCard - Сценарий обновления существующей банковской карты
type UpdateBankCard struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BankCardRepository - Реализация интерфейса BankCardRepositoryInterface
	BankCardRepository domain.BankCardRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u UpdateBankCard) Do(session domain.Session, cardID uuid.UUID, number, validThru, cvv, cardHolder string) error {
	card, err := u.BankCardRepository.Get(session.UserID, cardID)
	if err != nil {
		return err
	}

	if number != "" {
		card.Number = number
	}
	if validThru != "" {
		card.ValidThru = validThru
	}
	if cvv != "" {
		card.CVV = cvv
	}
	if cardHolder != "" {
		card.CardHolder = cardHolder
	}

	if err := u.Client.UpdateBankCard(session, &card); err != nil {
		return err
	}

	if err := u.BankCardRepository.Update(session.UserID, &card); err != nil {
		return err
	}

	return nil
}
