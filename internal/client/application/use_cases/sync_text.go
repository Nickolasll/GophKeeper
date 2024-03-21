package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// SyncText - Сценарий синхронизации текстовых данных
type SyncText struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// TextRepository - Реализация интерфейса TextRepositoryInterface
	TextRepository domain.TextRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u SyncText) Do(session domain.Session) error {
	texts, err := u.Client.GetAllTexts(session)
	if err != nil {
		return err
	}

	if err := u.TextRepository.ReplaceAll(session.UserID, texts); err != nil {
		return err
	}

	return nil
}
