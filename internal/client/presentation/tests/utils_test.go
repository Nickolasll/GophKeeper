// Package tests содержит интеграционные тесты клиента
package tests

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/urfave/cli/v3"
	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/application"
	"github.com/Nickolasll/goph-keeper/internal/client/config"
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
	binrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/binary_repository"
	jwkrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/jwk_repository"
	sessrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/session_repository"
	txtrepo "github.com/Nickolasll/goph-keeper/internal/client/infrastructure/test_repository"
	"github.com/Nickolasll/goph-keeper/internal/client/logger"
	"github.com/Nickolasll/goph-keeper/internal/client/presentation"
	"github.com/Nickolasll/goph-keeper/internal/crypto"
)

var db *bolt.DB
var cryptoService *crypto.CryptoService
var sessionRepository *sessrepo.SessionRepository
var textRepository *txtrepo.TextRepository
var jwkRepository *jwkrepo.JWKRepository
var binaryRepository *binrepo.BinaryRepository

func getJWKs() (jwk.Key, error) {
	jwks, err := jwk.FromRaw([]byte("My secret keys"))
	if err != nil {
		return jwks, err
	}

	return jwks, nil
}

func setup(client FakeHTTPClient) (*cli.Command, error) {
	var cmd *cli.Command
	var err error

	log := logger.New("./")

	cfg, err := config.New("./")
	if err != nil {
		return cmd, err
	}

	db, err = bolt.Open(cfg.DBFilePath, os.FileMode(cfg.DBFileMode), nil)
	if err != nil {
		return cmd, err
	}

	jwks, err := getJWKs()
	if err != nil {
		return cmd, err
	}

	publicKey, err := jwks.PublicKey()
	if err != nil {
		return cmd, err
	}

	certs, err := json.Marshal(publicKey)
	if err != nil {
		return cmd, err
	}

	cryptoService, err = crypto.New([]byte("1234567812345678"))
	if err != nil {
		return cmd, err
	}

	client.Certs = certs

	sessionRepository = sessrepo.New(db, cryptoService, log)
	textRepository = txtrepo.New(db, cryptoService, log)
	jwkRepository = jwkrepo.New(db, cryptoService, log)
	binaryRepository = binrepo.New(db, cryptoService, log)

	app := application.New(
		log,
		client,
		sessionRepository,
		textRepository,
		jwkRepository,
		binaryRepository,
	)

	cmd = presentation.New("v0.0.1", "01.01.1999", app, log, sessionRepository)

	return cmd, nil
}

func teardown() error {
	presentation.SetSession(nil)

	if err := sessionRepository.Delete(); err != nil {
		return err
	}

	if err := jwkRepository.Delete(); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func issueToken(userID uuid.UUID, exp time.Duration) (string, error) {
	jwks, err := getJWKs()
	if err != nil {
		return "", err
	}
	issuedAt := time.Now()
	expiration := issuedAt.Add(exp)
	token, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(expiration).
		Claim("UserID", userID.String()).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, jwks))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func createSession() (uuid.UUID, error) {
	userID := uuid.New()
	expiration := time.Hour
	token, err := issueToken(userID, expiration)
	if err != nil {
		return userID, err
	}

	session := domain.Session{
		UserID: userID,
		Token:  token,
	}
	err = sessionRepository.Save(session)
	if err != nil {
		return userID, err
	}
	presentation.SetSession(&session)

	return userID, nil
}

func createExpiredSession() error {
	userID := uuid.New()
	token, err := issueToken(userID, time.Duration(0))
	if err != nil {
		return err
	}

	session := domain.Session{
		UserID: userID,
		Token:  token,
	}
	err = sessionRepository.Save(session)
	if err != nil {
		return err
	}
	presentation.SetSession(&session)

	return nil
}
