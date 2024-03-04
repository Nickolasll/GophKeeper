// Package presentation содержит фабрику cli приложения и имплементацию команд
package presentation

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/Nickolasll/goph-keeper/internal/client/application"
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

var app *application.Application
var currentSession *domain.Session
var sessionRepository domain.SessionRepositoryInterface
var log *logrus.Logger

func New(
	_app *application.Application,
	_log *logrus.Logger,
	_sessionRepository domain.SessionRepositoryInterface,
) *cli.Command {
	var err error
	app = _app
	log = _log
	sessionRepository = _sessionRepository

	currentSession, err = sessionRepository.Get()
	if err != nil {
		currentSession = nil
	}

	cmdRegistration := registration()
	cmdLogin := login()

	cmdCreateText := createText()
	cmdUpdateText := updateText()
	cmdShowText := showText()

	cmdCreateBinary := createBinary()
	cmdUpdateBinary := updateBinary()
	cmdShowBinary := showBinary()

	cmd := cli.Command{
		Name:                  "gophkeeper",
		Version:               "v0.0.1",
		Usage:                 "Password and user data manager",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			&cmdRegistration,
			&cmdLogin,
			{
				Name:    "text",
				Aliases: []string{"t"},
				Usage:   "options for text",
				Commands: []*cli.Command{
					&cmdCreateText,
					&cmdUpdateText,
					&cmdShowText,
				},
			},
			{
				Name:    "binary",
				Aliases: []string{"b"},
				Usage:   "options for binary",
				Commands: []*cli.Command{
					&cmdCreateBinary,
					&cmdUpdateBinary,
					&cmdShowBinary,
				},
			},
		},
	}

	return &cmd
}

func SetSession(s *domain.Session) {
	currentSession = s
}
