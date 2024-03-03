// Package tests ...
package tests

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/application"
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
	"github.com/Nickolasll/goph-keeper/internal/client/infrastructure"
	"github.com/Nickolasll/goph-keeper/internal/client/presentation"
)

var db *bolt.DB
var sessionRepository infrastructure.SessionRepository
var textRepository infrastructure.TextRepository
var jwkRepository infrastructure.JWKRepository
var crypto application.CryptoService

type config struct {
	CryptoSecretKey []byte
	DBName          string
	DBFileMode      uint32
}

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
	cfg := config{
		CryptoSecretKey: []byte("1234567812345678"),
		DBName:          "test.db",
		DBFileMode:      0600,
	}

	log := logrus.New()

	db, err = bolt.Open(cfg.DBName, os.FileMode(cfg.DBFileMode), nil)
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

	crypto = application.CryptoService{
		SecretKey: cfg.CryptoSecretKey,
	}

	client.Certs = certs

	sessionRepository = infrastructure.SessionRepository{
		DB:     db,
		Crypto: crypto,
	}
	textRepository = infrastructure.TextRepository{
		DB:     db,
		Crypto: crypto,
	}
	jwkRepository = infrastructure.JWKRepository{
		DB:     db,
		Crypto: crypto,
	}

	app := application.CreateApplication(
		client,
		sessionRepository,
		textRepository,
		jwkRepository,
	)

	cmd = presentation.CLIFactory(&app, log, sessionRepository)

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

func issueToken(userID string, exp time.Duration) (string, error) {
	jwks, err := getJWKs()
	if err != nil {
		return "", err
	}
	issuedAt := time.Now()
	expiration := issuedAt.Add(exp)
	token, err := jwt.NewBuilder().IssuedAt(time.Now()).Expiration(expiration).Claim("UserID", userID).Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, jwks))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func createSession() (string, error) {
	userID := uuid.NewString()
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

func createExpiredSession() (string, error) {
	userID := uuid.NewString()
	token, err := issueToken(userID, time.Duration(0))
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
