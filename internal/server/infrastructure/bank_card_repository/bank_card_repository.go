// Package bankcardrepository содержит имлементацию интерфейса репозитория BankCardRepositoryInterface
package bankcardrepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// BankCardRepository - Имплементация репозитория для хранения банковских карт
type BankCardRepository struct {
	// DBPool - Интерфейс пула соединений pgxpool
	DBPool *pgxpool.Pool
	// Timeout - Таймаут операции
	Timeout time.Duration
	log     *logrus.Logger
}

// Create - Сохраняет новую банковскую карту
func (r BankCardRepository) Create(card *domain.BankCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		INSERT INTO bank_card_data
		(
			id
			, user_id
			, number
			, valid_thru
			, cvv
			, card_holder
		)
		VALUES
		(
			@id
			, @userID
			, @number
			, @valid_thru
			, @cvv
			, @card_holder
		)
		;`
	args := pgx.NamedArgs{
		"id":          card.ID,
		"userID":      card.UserID,
		"number":      card.Number,
		"valid_thru":  card.ValidThru,
		"cvv":         card.CVV,
		"card_holder": card.CardHolder,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Update - Обновляет существующую банковскую карту
func (r BankCardRepository) Update(card *domain.BankCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		UPDATE bank_card_data
		SET 
			number = @number
			, valid_thru = @valid_thru
			, cvv = @cvv
			, card_holder = @card_holder
		WHERE
			bank_card_data.id = @id
		    AND bank_card_data.user_id = @userID
		;`

	args := pgx.NamedArgs{
		"id":          card.ID,
		"userID":      card.UserID,
		"number":      card.Number,
		"valid_thru":  card.ValidThru,
		"cvv":         card.CVV,
		"card_holder": card.CardHolder,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Get - Возвращает банковскую карту по идентификатору пользователя и карты, если они существуют
func (r BankCardRepository) Get(userID, cardID uuid.UUID) (*domain.BankCard, error) {
	var card domain.BankCard

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			bank_card_data.id
			, bank_card_data.user_id
			, bank_card_data.number
			, bank_card_data.valid_thru
			, bank_card_data.cvv
			, bank_card_data.card_holder
		FROM
			bank_card_data
		WHERE
			bank_card_data.id = @cardID
		    AND bank_card_data.user_id = @userID
		;`
	args := pgx.NamedArgs{
		"cardID": cardID,
		"userID": userID,
	}
	err := r.DBPool.
		QueryRow(ctx, sql, args).
		Scan(&card.ID, &card.UserID, &card.Number, &card.ValidThru, &card.CVV, &card.CardHolder)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, err
	}

	return &card, err
}

// GetAll - Возвращает список банковских карт, принадлежащих пользователю
func (r BankCardRepository) GetAll(userID uuid.UUID) ([]*domain.BankCard, error) {
	result := []*domain.BankCard{}
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			bank_card_data.id
			, bank_card_data.user_id
			, bank_card_data.number
			, bank_card_data.valid_thru
			, bank_card_data.cvv
			, bank_card_data.card_holder
		FROM
			bank_card_data
		WHERE
			bank_card_data.user_id = @userID
		;`
	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := r.DBPool.Query(ctx, sql, args)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var card domain.BankCard
		err = rows.Scan(&card.ID, &card.UserID, &card.Number, &card.ValidThru, &card.CVV, &card.CardHolder)
		if err == nil {
			result = append(result, &card)
		}
	}
	if rows.Err() != nil {
		return result, err
	}

	return result, err
}

// New - Возвращает новый инстанс репозитория
func New(
	dbPool *pgxpool.Pool,
	timeout time.Duration,
	log *logrus.Logger,
) *BankCardRepository {
	return &BankCardRepository{
		DBPool:  dbPool,
		Timeout: timeout,
		log:     log,
	}
}
