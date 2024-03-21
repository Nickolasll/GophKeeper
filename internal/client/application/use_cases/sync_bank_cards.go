package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// SyncBankCards - Сценарий синхронизации банковских карт
type SyncBankCards struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BankCardRepository - Реализация интерфейса BankCardRepositoryInterface
	BankCardRepository domain.BankCardRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u SyncBankCards) Do(session domain.Session) error {
	cards, err := u.Client.GetAllBankCards(session)
	if err != nil {
		return err
	}

	if err := u.BankCardRepository.ReplaceAll(session.UserID, cards); err != nil {
		return err
	}

	return nil
}
