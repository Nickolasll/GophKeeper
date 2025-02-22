// Package main - инициализация сервера GophKeeper API
package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nickolasll/goph-keeper/internal/crypto"
	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
	"github.com/Nickolasll/goph-keeper/internal/server/config"
	bcardrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/bank_card_repository"
	binrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/binary_repository"
	crederepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/credentials_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/text_repository"
	usrrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/user_repository"
	"github.com/Nickolasll/goph-keeper/internal/server/logger"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"
)

// @Title GophKeeper API
// @Description Сервис хранения логинов, паролей, бинарных данных и прочей приватной информации.

// @version 0.0.1
// @BasePath /api/v1
// @Host 0.0.0.0:8080

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

// @Tag.name Status
// @Tag.description Группа запросов статуса сервера

// @Tag.name Auth
// @Tag.description Группа запросов авторизации

// @Tag.name Text
// @Tag.description Группа запросов для работы с текстовыми данными

// @Tag.name Binary
// @Tag.description Группа запросов для работы с бинарными данными

// @Tag.name Credentials
// @Tag.description Группа запросов для работы с логинами и паролями

// @Tag.name BankCard
// @Tag.description Группа запросов для работы с банковскими картами

// @Tag.name All
// @Tag.description Группа запросов для работы со всеми данными пользователя

func main() {
	log := logger.New()
	ctx := context.Background()

	cfg, err := config.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	joseService, err := jose.New(cfg.RawJWK, cfg.JWTExpiration, log)
	if err != nil {
		log.Fatal(err)
	}

	cryptoService, err := crypto.New(cfg.CryptoSecret)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepository := usrrepo.New(pool, cfg.DBTimeOut, log)
	textRepository := txtrepo.New(pool, cfg.DBTimeOut, log)
	binaryRepository := binrepo.New(pool, cfg.DBTimeOut, log)
	credentialsRepository := crederepo.New(pool, cfg.DBTimeOut, log)
	cardRepository := bcardrepo.New(pool, cfg.DBTimeOut, log)

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

	cert, err := tls.LoadX509KeyPair(cfg.X509CertPath, cfg.TLSKeyPath)
	if err != nil {
		log.Fatal(err) //nolint: gocritic
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		TLSConfig:         tlsConfig,
	}
	defer func() {
		err = server.Shutdown(ctx)
		if err != nil {
			log.Fatal(nil)
		}
	}()

	if err = server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
