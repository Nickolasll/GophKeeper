package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// ShowText - Сценарий получения всех локальных расшифрованных текстовых данных
type ShowText struct {
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// TextRepository - Реализация интерфейса TextRepositoryInterface
	TextRepository domain.TextRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u ShowText) Do(session domain.Session) ([]domain.Text, error) {
	result := []domain.Text{}
	_, err := u.CheckToken.Do(session.Token)
	if err != nil {
		return result, err
	}

	result, err = u.TextRepository.GetAll(session.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
