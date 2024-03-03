package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
	"github.com/Nickolasll/goph-keeper/internal/server/infrastructure"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Addr              string        `env:"ADDR, default=localhost:8080"`
	DBTimeOut         time.Duration `env:"DB_TIMEOUT, default=15s"`
	JWTExpiration     time.Duration `env:"JWT_EXPIRATION, default=600s"`
	RawJWK            []byte        `env:"RAW_JWK, default=My secret keys"`
	PostgresURL       string        `env:"POSTGRES_URL, default=postgresql://admin:admin@localhost:5432/gophkeeper"`
	CryptoSecret      []byte        `env:"CRYPTO_SECRET, default=1234567812345678"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT, default=600s"`
}

func main() {
	var cfg Config

	log := logrus.New()
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	key, err := jwk.FromRaw(cfg.RawJWK)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	jose := services.JOSEService{TokenExp: cfg.JWTExpiration, JWKS: key}
	crypto := services.CryptoService{
		SecretKey: cfg.CryptoSecret,
	}

	userRepository := infrastructure.UserRepository{
		DBPool:  pool,
		Timeout: cfg.DBTimeOut,
	}
	textRepository := infrastructure.TextRepository{
		DBPool:  pool,
		Timeout: cfg.DBTimeOut,
	}

	app := application.CreateApplication(
		userRepository,
		textRepository,
		jose,
		crypto,
	)

	router := presentation.ChiFactory(&app, &jose, log)

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
	if err = server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
