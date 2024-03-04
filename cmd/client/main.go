package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/application"
	"github.com/Nickolasll/goph-keeper/internal/client/infrastructure"
	"github.com/Nickolasll/goph-keeper/internal/client/presentation"
	"github.com/Nickolasll/goph-keeper/internal/crypto"
)

type config struct {
	CryptoSecretKey  []byte
	DBName           string
	DBFileMode       uint32
	ClientTimeoutSec int64
	ServerURL        string
}

func main() {
	var cmd *cli.Command
	var err error
	cfg := config{
		CryptoSecretKey:  []byte("1234567812345678"),
		DBName:           "user.db",
		DBFileMode:       0600,
		ClientTimeoutSec: 30,
		ServerURL:        "https://localhost:8080/api/v1/",
	}

	log := logrus.New()

	db, err := bolt.Open(cfg.DBName, os.FileMode(cfg.DBFileMode), nil)
	if err != nil {
		log.Fatal(err)
	}
	cert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
		RootCAs:       caCertPool,
		MinVersion:    tls.VersionTLS13,
	}

	timeout := time.Duration(cfg.ClientTimeoutSec) * time.Second
	client := infrastructure.HTTPClient{}.New(tlsConfig, timeout, cfg.ServerURL)

	cryptoService := crypto.New(cfg.CryptoSecretKey)

	sessionRepository := infrastructure.SessionRepository{
		DB:     db,
		Crypto: cryptoService,
	}
	textRepository := infrastructure.TextRepository{
		DB:     db,
		Crypto: cryptoService,
	}
	jwkRepository := infrastructure.JWKRepository{
		DB:     db,
		Crypto: cryptoService,
	}
	binaryRepository := infrastructure.BinaryRepository{
		DB:     db,
		Crypto: cryptoService,
	}

	app := application.New(
		client,
		sessionRepository,
		textRepository,
		jwkRepository,
		binaryRepository,
	)

	cmd = presentation.New(&app, log, sessionRepository)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
