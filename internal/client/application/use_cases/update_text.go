// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// UpdateText - Сценарий обновления существующих текстовых данных
type UpdateText struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// TextRepository - Реализация интерфейса TextRepository
	TextRepository domain.TextRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u UpdateText) Do(session domain.Session, textID uuid.UUID, content string) error {
	text, err := u.TextRepository.Get(session.UserID, textID)
	if err != nil {
		return err
	}

	text.Content = content

	if err := u.Client.UpdateText(session, text); err != nil {
		return err
	}

	if err := u.TextRepository.Update(session.UserID, text); err != nil {
		return err
	}

	return nil
}
