// Package tests ...
package tests

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nickolasll/goph-keeper/internal/crypto"
	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
	"github.com/Nickolasll/goph-keeper/internal/server/config"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
	bcardrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/bank_card_repository"
	binrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/binary_repository"
	crederepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/credentials_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/text_repository"
	usrrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/user_repository"
	"github.com/Nickolasll/goph-keeper/internal/server/logger"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"
)

var joseService *jose.JOSEService
var pool *pgxpool.Pool
var cryptoService *crypto.CryptoService
var userRepository *usrrepo.UserRepository
var textRepository *txtrepo.TextRepository
var binaryRepository *binrepo.BinaryRepository
var credentialsRepository *crederepo.CredentialsRepository
var cardRepository *bcardrepo.BankCardRepository

func setup() (*chi.Mux, error) {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	joseService, err = jose.New(cfg.RawJWK, cfg.JWTExpiration, log)
	if err != nil {
		return nil, err
	}

	pool, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	cryptoService, err = crypto.New(cfg.CryptoSecret)
	if err != nil {
		return nil, err
	}

	userRepository = usrrepo.New(pool, cfg.DBTimeOut, log)
	textRepository = txtrepo.New(pool, cfg.DBTimeOut, log)
	binaryRepository = binrepo.New(pool, cfg.DBTimeOut, log)
	credentialsRepository = crederepo.New(pool, cfg.DBTimeOut, log)
	cardRepository = bcardrepo.New(pool, cfg.DBTimeOut, log)

	app := application.New(
		log,
		joseService,
		cryptoService,
		userRepository,
		textRepository,
		binaryRepository,
		credentialsRepository,
		cardRepository,
	)

	router := presentation.New(app, joseService, log)

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
