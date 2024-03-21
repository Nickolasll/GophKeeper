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

// New - Фабрика CLI приложения
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
				Name:    "create",
				Usage:   "create text, binary, credentials or bank-cards",
				Aliases: []string{"c"},
				Commands: []*cli.Command{
					&cmdCreateText,
					&cmdCreateBinary,
					&cmdCreateCredentials,
					&cmdCreateBankCard,
				},
			},
			{
				Name:    "update",
				Usage:   "update text, binary, credentials or bank-cards",
				Aliases: []string{"u"},
				Commands: []*cli.Command{
					&cmdUpdateText,
					&cmdUpdateBinary,
					&cmdUpdateCredentials,
					&cmdUpdateBankCard,
				},
			},
			{
				Name:    "show",
				Usage:   "show local texts, binaries, credentials or bank-cards",
				Aliases: []string{"s"},
				Commands: []*cli.Command{
					&cmdShowText,
					&cmdShowBinary,
					&cmdShowCredentials,
					&cmdShowBankCard,
				},
			},
			{
				Name:  "sync",
				Usage: "manual override local data for text, binary, credentials or bank-cards",
				Commands: []*cli.Command{
					&cmdSyncText,
					&cmdSyncBinary,
					&cmdSyncCredentials,
					&cmdSyncBankCards,
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
