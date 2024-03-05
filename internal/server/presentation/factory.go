// Package presentation GophKeeper представляет собой серверную систему, позволяющую пользователю
// надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.
//
//		Schemes: https
//		Host: localhost
//		BasePath: /api/v1
//
//		components:
//			securitySchemes:
//		  		bearerAuth:
//		    		type: http
//		    		scheme: bearer
//		    		bearerFormat: JWT  # optional, for documentation purposes only
//
//		security:
//	  	- bearerAuth: []
//
// swagger:meta
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
	app = _app
	joseService = _jose
	log = _log

	validate = validator.New(validator.WithRequiredStructEnabled())

	router = chi.NewRouter()
	router.Use(logging)

	// swagger:route GET /health health
	// Return that server is alive
	//
	// responses:
	//  200: someResponse
	//  500: genericError
	router.Get("/api/v1/health", getHealthHandler)

	router.Post("/api/v1/auth/register", registrationHandler)
	router.Post("/api/v1/auth/login", loginHandler)
	router.Get("/api/v1/auth/certs", getCertsHandler)

	router.Post("/api/v1/text/create", auth(createTextHandler))
	router.Post("/api/v1/text/{textID}", auth(updateTextHandler))

	router.Post("/api/v1/binary/create", auth(createBinaryHandler))
	router.Post("/api/v1/binary/{binaryID}", auth(updateBinaryHandler))

	return router
}
