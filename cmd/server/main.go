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
	binrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/binary_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/text_repository"
	usrrepo "github.com/Nickolasll/goph-keeper/internal/server/infrastructure/user_repository"
	"github.com/Nickolasll/goph-keeper/internal/server/logger"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
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

	pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepository := usrrepo.New(pool, cfg.DBTimeOut, log)
	textRepository := txtrepo.New(pool, cfg.DBTimeOut, log)
	binaryRepository := binrepo.New(pool, cfg.DBTimeOut, log)

	app := application.New(
		log,
		joseService,
		cryptoService,
		userRepository,
		textRepository,
		binaryRepository,
	)

	router := presentation.New(app, joseService, log)

	cert, err := tls.LoadX509KeyPair(cfg.X509CertPath, cfg.X509KeyPath)
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
		err = server.Shutdown(context.Background())
		if err != nil {
			log.Fatal(nil)
		}
	}()

	if err = server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
