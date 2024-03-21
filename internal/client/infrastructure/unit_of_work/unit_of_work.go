// Package unitofwork обеспечивает транзакционность базы данных, изолируя слои и инкапсулируя всю логику
// работы с транзакциями
// На данный момент он работает только для транзакционного обновления всех пользовательских данных
// при синхронизации.
// Для исполнения других операций, необходимо переписать методы для записи в реализации репозиториев
package unitofwork

import (
	"errors"

	bolt "go.etcd.io/bbolt"

	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
	cardrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/bank_card_repository"
	binrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/binary_repository"
	credrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/credentials_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/text_repository"
)

// UnitOfWork - Интерфейс Unit Of Work для инкапсулирования транзакционной целостности
// По завершению работы транзакцию обязательно нужно коммитить или откатывать
type UnitOfWork struct {
	db                    *bolt.DB
	textRepository        txtrepo.TextRepository
	binaryRepository      binrepo.BinaryRepository
	credentialsRepository credrepo.CredentialsRepository
	bankCardRepository    cardrepo.BankCardRepository
	tx                    *bolt.Tx
	log                   *logrus.Logger
}

// Begin - Начало работы, создает транзакцию
func (uow *UnitOfWork) Begin() error {
	tx, err := uow.db.Begin(true)
	if err != nil {
		return err
	}

	uow.setTx(tx)

	return nil
}

func (uow *UnitOfWork) setTx(tx *bolt.Tx) {
	uow.tx = tx
	uow.textRepository.Tx = tx
	uow.binaryRepository.Tx = tx
	uow.credentialsRepository.Tx = tx
	uow.bankCardRepository.Tx = tx
}

// Commit - Выполняет коммит транзакции
func (uow *UnitOfWork) Commit() error {
	if uow.tx == nil {
		return bolt.ErrTxClosed
	}

	err := uow.tx.Commit()
	if err != nil {
		return err
	}

	uow.setTx(nil)

	return nil
}

// Rollback - Выполняет откат транзакции
func (uow *UnitOfWork) Rollback() error {
	if uow.tx == nil {
		return nil
	}
	err := uow.tx.Rollback()

	if err != nil && !errors.Is(err, bolt.ErrTxClosed) {
		return err
	}

	uow.setTx(nil)

	return nil
}

// TextRepository - Возвращает TextRepository для работы в пределах транзакции
func (uow *UnitOfWork) TextRepository() domain.TextRepositoryInterface {
	return uow.textRepository
}

// BinaryRepository - Возвращает BinaryRepository для работы в пределах транзакции
func (uow *UnitOfWork) BinaryRepository() domain.BinaryRepositoryInterface {
	return uow.binaryRepository
}

// CredentialsRepository - Возвращает CredentialsRepository для работы в пределах транзакции
func (uow *UnitOfWork) CredentialsRepository() domain.CredentialsRepositoryInterface {
	return uow.credentialsRepository
}

// BankCardRepository - Возвращает BankCardRepository для работы в пределах транзакции
func (uow *UnitOfWork) BankCardRepository() domain.BankCardRepositoryInterface {
	return uow.bankCardRepository
}

// New - Возвращает инстанс UnitOfWork
func New(
	db *bolt.DB,
	log *logrus.Logger,
	textRepository txtrepo.TextRepository,
	binaryRepository binrepo.BinaryRepository,
	credentialsRepository credrepo.CredentialsRepository,
	bankCardRepository cardrepo.BankCardRepository,
) *UnitOfWork {
	return &UnitOfWork{
		db:                    db,
		log:                   log,
		textRepository:        textRepository,
		binaryRepository:      binaryRepository,
		credentialsRepository: credentialsRepository,
		bankCardRepository:    bankCardRepository,
	}
}
