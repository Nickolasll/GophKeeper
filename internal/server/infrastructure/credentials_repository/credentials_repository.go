// Package credentialsrepository содержит имлементацию интерфейса репозитория CredentialsRepositoryInterface
package credentialsrepository

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

// CredentialsRepository - Имплементация репозитория для пар логин и пароль
type CredentialsRepository struct {
	// DBPool - Интерфейс пула соединений pgxpool
	DBPool *pgxpool.Pool
	// Timeout - Таймаут операции
	Timeout time.Duration
	log     *logrus.Logger
}

// Create - Сохраняет новую пару логин и пароль
func (r CredentialsRepository) Create(cred *domain.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		INSERT INTO credentials_data
		(
			id
			, user_id
			, name
			, login
			, password
			, meta
		)
		VALUES
		(
			@id
			, @userID
			, @name
			, @login
			, @password
			, @meta
		)
		;`
	args := pgx.NamedArgs{
		"id":       cred.ID,
		"userID":   cred.UserID,
		"name":     cred.Name,
		"login":    cred.Login,
		"password": cred.Password,
		"meta":     cred.Meta,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Update - Обновляет существующую пару логин и пароль
func (r CredentialsRepository) Update(cred *domain.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		UPDATE credentials_data
		SET 
			name = @name
			, login = @login
			, password = @password
			, meta = @meta
		WHERE
			credentials_data.id = @id
		    AND credentials_data.user_id = @userID
		;`

	args := pgx.NamedArgs{
		"id":       cred.ID,
		"userID":   cred.UserID,
		"name":     cred.Name,
		"login":    cred.Login,
		"password": cred.Password,
		"meta":     cred.Meta,
	}
	_, err := r.DBPool.Exec(ctx, sql, args)

	return err
}

// Get - Возвращает пару логин и пароль по идентификатору пользователя и данных, если они существуют
func (r CredentialsRepository) Get(userID, credID uuid.UUID) (*domain.Credentials, error) {
	var cred domain.Credentials

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			credentials_data.id
			, credentials_data.user_id
			, credentials_data.name
			, credentials_data.login
			, credentials_data.password
			, credentials_data.meta
		FROM
			credentials_data
		WHERE
			credentials_data.id = @credID
		    AND credentials_data.user_id = @userID
		;`
	args := pgx.NamedArgs{
		"credID": credID,
		"userID": userID,
	}
	err := r.DBPool.
		QueryRow(ctx, sql, args).
		Scan(
			&cred.ID,
			&cred.UserID,
			&cred.Name,
			&cred.Login,
			&cred.Password,
			&cred.Meta,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, err
	}

	return &cred, err
}

// GetAll - Возвращает список пар логин и пароль, принадлежащих пользователю
func (r CredentialsRepository) GetAll(userID uuid.UUID) ([]*domain.Credentials, error) {
	result := []*domain.Credentials{}
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()
	sql := `
		SELECT
			credentials_data.id
			, credentials_data.user_id
			, credentials_data.name
			, credentials_data.login
			, credentials_data.password
			, credentials_data.meta
		FROM
			credentials_data
		WHERE
			credentials_data.user_id = @userID
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
		var cred domain.Credentials
		err = rows.Scan(
			&cred.ID,
			&cred.UserID,
			&cred.Name,
			&cred.Login,
			&cred.Password,
			&cred.Meta,
		)
		if err == nil {
			result = append(result, &cred)
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
) *CredentialsRepository {
	return &CredentialsRepository{
		DBPool:  dbPool,
		Timeout: timeout,
		log:     log,
	}
}
