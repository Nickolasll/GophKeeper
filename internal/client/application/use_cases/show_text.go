// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// ShowText - Сценарий получения расшифрованных текстовых данных
type ShowText struct {
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// TextRepository - Реализация интерфейса TextRepository
	TextRepository domain.TextRepositoryInterface
}

// Execute - Вызов логики сценария использования
func (u ShowText) Execute(session domain.Session) ([]domain.Text, error) {
	result := []domain.Text{}
	_, err := u.CheckToken.Execute(session.Token)
	if err != nil {
		return result, err
	}

	result, err = u.TextRepository.GetAll(session.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
