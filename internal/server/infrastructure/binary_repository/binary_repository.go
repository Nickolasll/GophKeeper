// Package binaryrepository содержит имлементацию интерфейса репозитория BinaryRepositoryInterface
package binaryrepository

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

// BinaryRepository - Имплементация репозитория для произвольных бинарных данных
type BinaryRepository struct {
	// DBPool - Интерфейс пула соединений pgxpool
	DBPool *pgxpool.Pool
	// Timeout - Таймаут операции
	Timeout time.Duration
	log     *logrus.Logger
}

// Create - Сохраняет новые бинарные данных
func (r BinaryRepository) Create(bin domain.Binary) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		INSERT INTO binary_data
		(
			id
			, user_id
			, content
		)
		VALUES
		(
			@id
			, @userID
			, @content
		)
		;`
	args := pgx.NamedArgs{
		"id":      bin.ID,
		"userID":  bin.UserID,
		"content": bin.Content,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Update - Обновляет существующие бинарные данные
func (r BinaryRepository) Update(bin domain.Binary) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		UPDATE binary_data
		SET 
			content = @content
		WHERE
		    binary_data.id = @id
		    AND binary_data.user_id = @userID
		;`

	args := pgx.NamedArgs{
		"id":      bin.ID,
		"userID":  bin.UserID,
		"content": bin.Content,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Get - Возвращает бинарные данные по идентификатору пользователя и данных, если они существуют
func (r BinaryRepository) Get(userID, binID uuid.UUID) (*domain.Binary, error) {
	var bin domain.Binary

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
		    binary_data.id
			, binary_data.user_id
			, binary_data.content
		FROM
		    binary_data
		WHERE
		    binary_data.id = @binID
		    AND binary_data.user_id = @userID
		;`
	args := pgx.NamedArgs{
		"binID":  binID,
		"userID": userID,
	}
	err := r.DBPool.QueryRow(ctx, sql, args).Scan(&bin.ID, &bin.UserID, &bin.Content)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, err
	}

	return &bin, err
}

// GetAll - Возвращает список бинарных данных, принадлежащих пользователю
func (r BinaryRepository) GetAll(userID uuid.UUID) ([]domain.Binary, error) {
	result := []domain.Binary{}
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
		    binary_data.id
			, binary_data.user_id
			, binary_data.content
		FROM
		    binary_data
		WHERE
		    binary_data.user_id = @userID
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
		var bin domain.Binary
		err = rows.Scan(&bin)
		if err == nil {
			result = append(result, bin)
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
) *BinaryRepository {
	return &BinaryRepository{
		DBPool:  dbPool,
		Timeout: timeout,
		log:     log,
	}
}
