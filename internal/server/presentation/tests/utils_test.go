// Package tests ...
package tests

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
	"github.com/Nickolasll/goph-keeper/internal/server/infrastructure"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"
)

var jose services.JOSEService
var pool *pgxpool.Pool
var crypto services.CryptoService
var userRepository infrastructure.UserRepository
var textRepository infrastructure.TextRepository
var binaryRepository infrastructure.BinaryRepository

type config struct {
	TimeoutDuration time.Duration
	CryptoSecretKey []byte
	RawJWK          []byte
	PostgresURL     string
}

func setup() (*chi.Mux, error) {
	cfg := config{
		TimeoutDuration: time.Duration(30) * time.Second,
		CryptoSecretKey: []byte("1234567812345678"),
		RawJWK:          []byte("My secret keys"),
		PostgresURL:     "postgresql://admin:admin@localhost:5432/gophkeeper",
	}

	key, err := jwk.FromRaw(cfg.RawJWK)
	if err != nil {
		return nil, err
	}
	jose = services.JOSEService{
		TokenExp: cfg.TimeoutDuration,
		JWKS:     key,
	}

	pool, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	crypto = services.CryptoService{
		SecretKey: cfg.CryptoSecretKey,
	}

	userRepository = infrastructure.UserRepository{
		DBPool:  pool,
		Timeout: cfg.TimeoutDuration,
	}
	textRepository = infrastructure.TextRepository{
		DBPool:  pool,
		Timeout: cfg.TimeoutDuration,
	}
	binaryRepository = infrastructure.BinaryRepository{
		DBPool:  pool,
		Timeout: cfg.TimeoutDuration,
	}

	app := application.CreateApplication(
		jose,
		crypto,
		userRepository,
		textRepository,
		binaryRepository,
	)
	log := logrus.New()
	router := presentation.ChiFactory(&app, &jose, log)

	return router, nil
}

func teardown() {
	pool.Close()
}

func createUser(id uuid.UUID) error {
	user := domain.User{
		ID:       id,
		Login:    uuid.NewString(),
		Password: "password",
	}

	return userRepository.Create(user)
}
