// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// ShowBinary - Сценарий получения расшифрованных бинарных данных
type ShowBinary struct {
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// BinaryRepository - Реализация интерфейса BinaryRepositoryInterface
	BinaryRepository domain.BinaryRepositoryInterface
}

// Execute - Вызов логики сценария использования
func (u ShowBinary) Execute(session domain.Session) ([]domain.Binary, error) {
	result := []domain.Binary{}
	_, err := u.CheckToken.Execute(session.Token)
	if err != nil {
		return result, err
	}

	result, err = u.BinaryRepository.GetAll(session.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
