package main

import (
	_ "embed"
	"path/filepath"

	"context"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/application"
	"github.com/Nickolasll/goph-keeper/internal/client/config"
	binrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/binary_repository"
	httpclient "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/http_client"
	jwkrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/jwk_repository"
	sessrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/session_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/test_repository"
	"github.com/Nickolasll/goph-keeper/internal/client/logger"
	"github.com/Nickolasll/goph-keeper/internal/client/presentation"
	"github.com/Nickolasll/goph-keeper/internal/crypto"
)

//go:embed ca.crt
var caCRT []byte

//go:embed secret.key
var secret []byte

var (
	Version   string
	BuildDate string
)

func main() {
	ex, _ := os.Executable()
	root := filepath.Dir(ex)

	log := logger.New(root)

	cfg, err := config.New(root)
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(cfg.DBFilePath, os.FileMode(cfg.DBFileMode), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := httpclient.New(
		log,
		caCRT,
		cfg.ClientTimeout,
		cfg.ServerURL+cfg.APIVer,
	)

	cryptoService, err := crypto.New(secret)
	if err != nil {
		log.Fatal(err) //nolint: gocritic
	}

	sessionRepository := sessrepo.New(db, cryptoService, log)
	textRepository := txtrepo.New(db, cryptoService, log)
	jwkRepository := jwkrepo.New(db, cryptoService, log)
	binaryRepository := binrepo.New(db, cryptoService, log)

	app := application.New(
		log,
		client,
		sessionRepository,
		textRepository,
		jwkRepository,
		binaryRepository,
	)

	cmd := presentation.New(Version, BuildDate, app, log, sessionRepository)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
