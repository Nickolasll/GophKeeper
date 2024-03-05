// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CreateText - Сценарий создания новых текстовых данных
type CreateText struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// TextRepository - Реализация интерфейса TextRepository
	TextRepository domain.TextRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u CreateText) Do(session domain.Session, content string) error {
	textID, err := u.Client.CreateText(session, content)
	if err != nil {
		return err
	}

	text := domain.Text{
		ID:      textID,
		Content: content,
	}

	if err := u.TextRepository.Create(session.UserID, text); err != nil {
		return err
	}

	return nil
}
