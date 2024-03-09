package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// GetAll - Сценарий использования для получения
type GetAll struct {
	// GetAllTexts - Сценарий использования для получения всех расшифрованных текстовых данных
	GetAllTexts GetAllTexts
	// GetAllBankCards - Сценарий использования для получения всех расшифрованных банковских карт
	GetAllBankCards GetAllBankCards
	// GetAllBinaries - Сценарий использования для получения всех расшифрованных бинарных данных
	GetAllBinaries GetAllBinaries
	// GetAllCredentials - Сценарий использования для получения всех расшифрованных логинов и паролей
	GetAllCredentials GetAllCredentials
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает все расшифрованные данные пользователя
func (u *GetAll) Do(userID uuid.UUID) (
	texts []domain.Text,
	bankCards []*domain.BankCard,
	binaries []domain.Binary,
	credentials []domain.Credentials,
	err error,
) {
	g := errgroup.Group{}
	g.Go(func() error {
		res, err := u.GetAllTexts.Do(userID)
		if err != nil {
			return err
		}
		texts = res

		return nil
	})
	g.Go(func() error {
		res, err := u.GetAllBinaries.Do(userID)
		if err != nil {
			return err
		}
		binaries = res

		return nil
	})
	g.Go(func() error {
		res, err := u.GetAllCredentials.Do(userID)
		if err != nil {
			return err
		}
		credentials = res

		return nil
	})
	g.Go(func() error {
		res, err := u.GetAllBankCards.Do(userID)
		if err != nil {
			return err
		}
		bankCards = res

		return nil
	})
	err = g.Wait()
	if err != nil {
		return []domain.Text{}, []*domain.BankCard{}, []domain.Binary{}, []domain.Credentials{}, err
	}

	return texts, bankCards, binaries, credentials, nil
}
