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
	version string,
	buildDate string,
	_app *application.Application,
	_log *logrus.Logger,
	_sessionRepository domain.SessionRepositoryInterface,
) *cli.Command {
	var err error
	app = _app
	log = _log
	sessionRepository = _sessionRepository

	regexpMustCompile()

	currentSession, err = sessionRepository.Get()
	if err != nil {
		currentSession = nil
	}

	cmdRegistration := registration()
	cmdLogin := login()

	cmdCreateText := createText()
	cmdUpdateText := updateText()
	cmdShowText := showText()
	cmdSyncText := syncText()

	cmdCreateBinary := createBinary()
	cmdUpdateBinary := updateBinary()
	cmdShowBinary := showBinary()
	cmdSyncBinary := syncBinary()

	cmdCreateCredentials := createCredentials()
	cmdUpdateCredentials := updateCredentials()
	cmdShowCredentials := showCredentials()
	cmdSyncCredentials := syncCredentials()

	cmdCreateBankCard := createBankCard()
	cmdUpdateBankCard := updateBankCard()
	cmdShowBankCard := showBankCards()
	cmdSyncBankCards := syncBankCards()

	cmdSyncAll := syncAll()

	cmd := cli.Command{
		Name:                  "gophkeeper",
		Version:               version + ", build at: " + buildDate,
		Usage:                 "Password and user data manager",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			&cmdRegistration,
			&cmdLogin,
			{
				Name:  "text",
				Usage: "options for text",
				Commands: []*cli.Command{
					&cmdCreateText,
					&cmdUpdateText,
					&cmdShowText,
					&cmdSyncText,
				},
			},
			{
				Name:  "binary",
				Usage: "options for binary",
				Commands: []*cli.Command{
					&cmdCreateBinary,
					&cmdUpdateBinary,
					&cmdShowBinary,
					&cmdSyncBinary,
				},
			},
			{
				Name:  "credentials",
				Usage: "options for credentials",
				Commands: []*cli.Command{
					&cmdCreateCredentials,
					&cmdUpdateCredentials,
					&cmdShowCredentials,
					&cmdSyncCredentials,
				},
			},
			{
				Name:  "bank_card",
				Usage: "options for bank cards",
				Commands: []*cli.Command{
					&cmdCreateBankCard,
					&cmdUpdateBankCard,
					&cmdShowBankCard,
					&cmdSyncBankCards,
				},
			},
			{
				Name:  "all",
				Usage: "options for all",
				Commands: []*cli.Command{
					&cmdSyncAll,
				},
			},
		},
	}

	return &cmd
}

func SetSession(s *domain.Session) {
	currentSession = s
}
