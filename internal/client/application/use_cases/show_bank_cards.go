// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// ShowBankCards - Сценарий получения расшифрованных банковских карт
type ShowBankCards struct {
	// CheckToken - Сценарий проверки JWT, возвращает UserID в формате строки
	CheckToken *CheckToken
	// BankCardRepository - Реализация интерфейса BankCardRepository
	BankCardRepository domain.BankCardRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u ShowBankCards) Do(session domain.Session) ([]domain.BankCard, error) {
	result := []domain.BankCard{}
	_, err := u.CheckToken.Do(session.Token)
	if err != nil {
		return result, err
	}

	result, err = u.BankCardRepository.GetAll(session.UserID)
	if err != nil {
		return result, err
	}

	return result, nil
}
