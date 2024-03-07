// Package presentation содержит фабрику роутера и реализацию обработчиков запросов
package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
)

var app *application.Application
var log *logrus.Logger
var joseService *jose.JOSEService
var validate *validator.Validate
var router *chi.Mux

// New - Фабрика HTTP роутера
func New(
	_app *application.Application,
	_jose *jose.JOSEService,
	_log *logrus.Logger,
) *chi.Mux {
	var err error

	app = _app
	joseService = _jose
	log = _log

	validate, err = newValidator()
	if err != nil {
		log.Fatal(err)
	}

	router = chi.NewRouter()
	router.Use(logging)

	router.Get("/api/v1/health", getHealthHandler)

	router.Post("/api/v1/auth/register", registrationHandler)
	router.Post("/api/v1/auth/login", loginHandler)
	router.Get("/api/v1/auth/certs", getCertsHandler)

	router.Post("/api/v1/text/create", auth(createTextHandler))
	router.Post("/api/v1/text/{textID}", auth(updateTextHandler))
	router.Get("/api/v1/text/all", auth(getAllTextsHandler))

	router.Post("/api/v1/binary/create", auth(createBinaryHandler))
	router.Post("/api/v1/binary/{binaryID}", auth(updateBinaryHandler))

	router.Post("/api/v1/credentials/create", auth(createCredentialsHandler))
	router.Post("/api/v1/credentials/{credID}", auth(updateCredentialsHandler))

	router.Post("/api/v1/bank_card/create", auth(createBankCardHandler))
	router.Post("/api/v1/bank_card/{cardID}", auth(updateBankCardHandler))

	return router
}
