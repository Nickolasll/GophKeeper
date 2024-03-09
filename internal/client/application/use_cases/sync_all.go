package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// SyncAll - Сценарий синхронизации всех пользовательских данных
type SyncAll struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// UnitOfWork - Реализация интерфейса UnitOfWorkInterface
	UnitOfWork domain.UnitOfWorkInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u SyncAll) Do(session domain.Session) error {
	texts, bankCards, binaries, credentials, err := u.Client.GetAll(session)
	if err != nil {
		return err
	}

	err = u.UnitOfWork.Begin()
	if err != nil {
		return err
	}
	defer u.UnitOfWork.Rollback() // nolint: errcheck

	err = u.UnitOfWork.TextRepository().ReplaceAll(session.UserID, texts)
	if err != nil {
		return err
	}

	err = u.UnitOfWork.BankCardRepository().ReplaceAll(session.UserID, bankCards)
	if err != nil {
		return err
	}

	err = u.UnitOfWork.BinaryRepository().ReplaceAll(session.UserID, binaries)
	if err != nil {
		return err
	}

	err = u.UnitOfWork.CredentialsRepository().ReplaceAll(session.UserID, credentials)
	if err != nil {
		return err
	}

	err = u.UnitOfWork.Commit()
	if err != nil {
		return err
	}

	return nil
}
