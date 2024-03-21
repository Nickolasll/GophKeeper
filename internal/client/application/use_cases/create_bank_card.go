package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CreateBankCard - Сценарий создания новой банковской карты
type CreateBankCard struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BankCardRepository - Реализация интерфейса BankCardRepositoryInterface
	BankCardRepository domain.BankCardRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u CreateBankCard) Do(
	session domain.Session,
	number, validThru, cvv, cardHolder, meta string,
) error {
	cardID, err := u.Client.CreateBankCard(session, number, validThru, cvv, cardHolder, meta)
	if err != nil {
		return err
	}

	card := domain.BankCard{
		ID:         cardID,
		Number:     number,
		ValidThru:  validThru,
		CVV:        cvv,
		CardHolder: cardHolder,
		Meta:       meta,
	}

	if err := u.BankCardRepository.Create(session.UserID, &card); err != nil {
		return err
	}

	return nil
}
