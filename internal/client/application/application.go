// Package application содержит фабрику приложения
package application

import (
	"github.com/sirupsen/logrus"

	usecases "github.com/Nickolasll/goph-keeper/internal/client/application/use_cases"
	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// Application - Приложение, инкапсулирует всю доступную бизнес логику
type Application struct {
	// Registration - Сценарий регистрации по логину и паролю
	Registration usecases.Registration
	// Login - Сценарий входа по логину и паролю
	Login usecases.Login
	// CreateText - Сценарий создания новых текстовых данных
	CreateText usecases.CreateText
	// UpdateText - Сценарий обновления существующих текстовых данных
	UpdateText usecases.UpdateText
	// ShowText - Сценарий получения расшифрованных текстовых данных
	ShowText usecases.ShowText
	// Перезапись текущих пользовательских данных
	SyncText usecases.SyncText
	// CreateBinary - Сценарий создания новых бинарных данных
	CreateBinary usecases.CreateBinary
	// UpdateBinary - Сценарий обновления существующих бинарных данных
	UpdateBinary usecases.UpdateBinary
	// ShowBinary - Сценарий получения расшифрованных бинарных данных
	ShowBinary usecases.ShowBinary
	// CreateCredentials - Сценарий создания новой пары логин и пароль
	CreateCredentials usecases.CreateCredentials
	// UpdateCredentials - Сценарий обновления существующей пары логин и пароль
	UpdateCredentials usecases.UpdateCredentials
	// ShowCredentials - Сценарий получения расшифрованных логинов и паролей пользователя
	ShowCredentials usecases.ShowCredentials
	// CreateBankCard - Сценарий создания новой банковской карты
	CreateBankCard usecases.CreateBankCard
	// UpdateCredentials - Сценарий обновления существующей банковской карты
	UpdateBankCard usecases.UpdateBankCard
	// ShowBankCards - Сценарий получения расшифрованных банковских карт
	ShowBankCards usecases.ShowBankCards
}

// New - Фабрика приложения
func New(
	log *logrus.Logger,
	client domain.GophKeeperClientInterface,
	sessionRepository domain.SessionRepositoryInterface,
	textRepository domain.TextRepositoryInterface,
	jwkRepository domain.JWKRepositoryInterface,
	binaryRepository domain.BinaryRepositoryInterface,
	credentialsRepository domain.CredentialsRepositoryInterface,
	bankCardRepository domain.BankCardRepositoryInterface,
) *Application {
	jwk, err := jwkRepository.Get()

	if err != nil {
		// Если не нашли ключа в репозитории, клиент запросит его с сервера
		jwk = nil
	}

	checkToken := usecases.CheckToken{
		Client:        client,
		JWKRepository: jwkRepository,
		Key:           jwk,
		Log:           log,
	}
	registration := usecases.Registration{
		Client:            client,
		SessionRepository: sessionRepository,
		CheckToken:        &checkToken,
		Log:               log,
	}
	login := usecases.Login{
		Client:            client,
		SessionRepository: sessionRepository,
		CheckToken:        &checkToken,
		Log:               log,
	}

	createText := usecases.CreateText{
		Client:         client,
		TextRepository: textRepository,
		Log:            log,
	}
	updateText := usecases.UpdateText{
		Client:         client,
		TextRepository: textRepository,
		Log:            log,
	}
	showText := usecases.ShowText{
		CheckToken:     &checkToken,
		TextRepository: textRepository,
		Log:            log,
	}
	syncText := usecases.SyncText{
		Client:         client,
		TextRepository: textRepository,
		Log:            log,
	}

	createBinary := usecases.CreateBinary{
		Client:           client,
		BinaryRepository: binaryRepository,
		Log:              log,
	}
	updateBinary := usecases.UpdateBinary{
		Client:           client,
		BinaryRepository: binaryRepository,
		Log:              log,
	}
	showBinary := usecases.ShowBinary{
		CheckToken:       &checkToken,
		BinaryRepository: binaryRepository,
		Log:              log,
	}

	createCredentials := usecases.CreateCredentials{
		Client:                client,
		CredentialsRepository: credentialsRepository,
		Log:                   log,
	}
	updateCredentials := usecases.UpdateCredentials{
		Client:                client,
		CredentialsRepository: credentialsRepository,
		Log:                   log,
	}
	showCredentials := usecases.ShowCredentials{
		CheckToken:            &checkToken,
		CredentialsRepository: credentialsRepository,
		Log:                   log,
	}

	createBankCard := usecases.CreateBankCard{
		Client:             client,
		BankCardRepository: bankCardRepository,
		Log:                log,
	}
	updateBankCard := usecases.UpdateBankCard{
		Client:             client,
		BankCardRepository: bankCardRepository,
		Log:                log,
	}
	showBankCards := usecases.ShowBankCards{
		CheckToken:         &checkToken,
		BankCardRepository: bankCardRepository,
		Log:                log,
	}

	return &Application{
		Registration:      registration,
		Login:             login,
		CreateText:        createText,
		UpdateText:        updateText,
		ShowText:          showText,
		SyncText:          syncText,
		CreateBinary:      createBinary,
		UpdateBinary:      updateBinary,
		ShowBinary:        showBinary,
		CreateCredentials: createCredentials,
		UpdateCredentials: updateCredentials,
		ShowCredentials:   showCredentials,
		CreateBankCard:    createBankCard,
		UpdateBankCard:    updateBankCard,
		ShowBankCards:     showBankCards,
	}
}
