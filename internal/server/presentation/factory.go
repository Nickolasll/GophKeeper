package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application"
	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
)

var app *application.Application
var log *logrus.Logger
var jose *services.JOSEService
var validate *validator.Validate
var router *chi.Mux

// ChiFactory - Фабрика HTTP роутера
func ChiFactory(
	_app *application.Application,
	_jose *services.JOSEService,
	_log *logrus.Logger,
) *chi.Mux {
	app = _app
	jose = _jose
	log = _log

	validate = validator.New(validator.WithRequiredStructEnabled())

	router = chi.NewRouter()
	router.Use(logging)

	router.Post("/api/v1/auth/register", registrationHandler)
	router.Post("/api/v1/auth/login", loginHandler)
	router.Get("/api/v1/auth/certs", getCertsHandler)
	router.Post("/api/v1/text/create", auth(createTextHandler))
	router.Post("/api/v1/text/{textID}", auth(updateTextHandler))

	return router
}
