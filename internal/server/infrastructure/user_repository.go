// Package infrastructure содержит имлементацию репозиториев
package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UserRepository - Имлементация репозитория текстовых данных
type UserRepository struct {
	// DBPool - Пул соединений pgx
	DBPool *pgxpool.Pool
	// Timeout - Таймаут операции
	Timeout time.Duration
}

// Create - Сохраняет нового пользователя
func (r UserRepository) Create(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		INSERT INTO users
		(
			id
			, login
			, password
		)
		VALUES 
		(
			@id
			, @login
			, @password
		)
		;`
	args := pgx.NamedArgs{
		"id":       user.ID,
		"login":    user.Login,
		"password": user.Password,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// GetByLogin - Возвращает пользователя по логину, если он существует
func (r UserRepository) GetByLogin(login string) (*domain.User, error) {
	var user domain.User
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			users.id
			, users.login
			, users.password
		FROM
			users
		WHERE
			users.login = @login
		;`
	args := pgx.NamedArgs{
		"login": login,
	}
	err := r.DBPool.QueryRow(ctx, sql, args).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, err
	}

	return &user, nil
}
