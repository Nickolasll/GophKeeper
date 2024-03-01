// Package infrastructure содержит имлементацию репозиториев
package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// TextRepository - Имплементация репозитория для произвольных текстовых данных
type TextRepository struct {
	// DBPool - Интерфейс пула соединений pgxpool
	DBPool *pgxpool.Pool
	// Timeout - Таймаут операции
	Timeout time.Duration
}

// Create - Сохраняет новые текстовые данные
func (r TextRepository) Create(text domain.Text) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		INSERT INTO text
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
		"id":      text.ID,
		"userID":  text.UserID,
		"content": text.Content,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Update - Обновляет текстовые данные
func (r TextRepository) Update(text domain.Text) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		UPDATE text
		SET 
			content = @content
		WHERE
			text.id = @id
		    AND text.user_id = @userID
		;`

	args := pgx.NamedArgs{
		"id":      text.ID,
		"userID":  text.UserID,
		"content": text.Content,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Get - Возвращает текстовые данные по идентификатору данных и пользователя, если они существуют
func (r TextRepository) Get(textID, userID uuid.UUID) (*domain.Text, error) {
	var text domain.Text

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			text.id
			, text.user_id
			, text.content
		FROM
			text
		WHERE
			text.id = @textID
		    AND text.user_id = @userID
		;`
	args := pgx.NamedArgs{
		"textID": textID,
		"userID": userID,
	}
	err := r.DBPool.QueryRow(ctx, sql, args).Scan(&text.ID, &text.UserID, &text.Content)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, err
	}

	return &text, err
}

// FindByUserID - Возвращает список текстовых данных, принадлежащих пользователю
func (r TextRepository) FindByUserID(userID uuid.UUID) ([]domain.Text, error) {
	result := []domain.Text{}
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			text.id
			, text.user_id
			, text.content
		FROM
			text
		WHERE
			text.user_id = @userID
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
		var text domain.Text
		err = rows.Scan(&text)
		if err == nil {
			result = append(result, text)
		}
	}
	if rows.Err() != nil {
		return result, err
	}

	return result, err
}
