// Package application содержит фабрику приложения
package application

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/application/jose"
	usecases "github.com/Nickolasll/goph-keeper/internal/server/application/use_cases"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// Application - Приложение, инкапсулирует всю доступную бизнес логику
type Application struct {
	// Registration - Сценарий использования регистрации пользователя по логину и паролю
	Registration usecases.Registration
	// Login - Сценарий использования входа пользователя по логину и паролю
	Login usecases.Login
	// CreateText - Сценарий использования для создания зашифрованных текстовых данных
	CreateText usecases.CreateText
	// UpdateText - Сценарий использования для обновления существующих зашифрованных текстовых данных
	UpdateText usecases.UpdateText
	// GetAllTexts - Получение всех расшифрованных текстовых данных
	GetAllTexts usecases.GetAllTexts
	// CreateText - Сценарий использования для создания зашифрованных бинарных данных
	CreateBinary usecases.CreateBinary
	// UpdateText - Сценарий использования для обновления существующих зашифрованных бинарных данных
	UpdateBinary usecases.UpdateBinary
	// GetAllBinaries - Получение всех расшифрованных бинарных данных
	GetAllBinaries usecases.GetAllBinaries
	// CreateCredentials - Сценарий использования для создания зашифрованной пары логин и пароль
	CreateCredentials usecases.CreateCredentials
	// UpdateCredentials - Сценарий использования для обновления существующей зашифрованной пары логин и пароль
	UpdateCredentials usecases.UpdateCredentials
	// GetAllCredentials - Получение всех расшифрованных логинов и паролей
	GetAllCredentials usecases.GetAllCredentials
	// CreateBankCard - Сценарий использования для создания зашифрованной банковской карты
	CreateBankCard usecases.CreateBankCard
	// UpdateBankCard - Сценарий использования для обновления существующей зашифрованной банковской карты
	UpdateBankCard usecases.UpdateBankCard
	// GetAllBankCards - Получение всех расшифрованных банковских карт
	GetAllBankCards usecases.GetAllBankCards
	GetAll          usecases.GetAll
}

// New - Фабрика приложения
func New(
	log *logrus.Logger,
	joseService *jose.JOSEService,
	crypto domain.CryptoServiceInterface,
	userRepository domain.UserRepositoryInterface,
	textRepository domain.TextRepositoryInterface,
	binaryRepository domain.BinaryRepositoryInterface,
	credentialsRepository domain.CredentialsRepositoryInterface,
	bankCardRepository domain.BankCardRepositoryInterface,
) *Application {
	registration := usecases.Registration{
		UserRepository: userRepository,
		JOSE:           joseService,
		Log:            log,
	}

	login := usecases.Login{
		UserRepository: userRepository,
		JOSE:           joseService,
		Log:            log,
	}

	createText := usecases.CreateText{
		TextRepository: textRepository,
		Crypto:         crypto,
		Log:            log,
	}
	updateText := usecases.UpdateText{
		TextRepository: textRepository,
		Crypto:         crypto,
		Log:            log,
	}
	getAllTexts := usecases.GetAllTexts{
		TextRepository: textRepository,
		Crypto:         crypto,
		Log:            log,
	}

	createBinary := usecases.CreateBinary{
		BinaryRepository: binaryRepository,
		Crypto:           crypto,
		Log:              log,
	}
	updateBinary := usecases.UpdateBinary{
		BinaryRepository: binaryRepository,
		Crypto:           crypto,
		Log:              log,
	}
	getAllBinaries := usecases.GetAllBinaries{
		BinaryRepository: binaryRepository,
		Crypto:           crypto,
		Log:              log,
	}

	createCredentials := usecases.CreateCredentials{
		CredentialsRepository: credentialsRepository,
		Crypto:                crypto,
		Log:                   log,
	}
	updateCredentials := usecases.UpdateCredentials{
		CredentialsRepository: credentialsRepository,
		Crypto:                crypto,
		Log:                   log,
	}
	getAllCredentials := usecases.GetAllCredentials{
		CredentialsRepository: credentialsRepository,
		Crypto:                crypto,
		Log:                   log,
	}

	createBankCard := usecases.CreateBankCard{
		BankCardRepository: bankCardRepository,
		Crypto:             crypto,
		Log:                log,
	}
	updateBankCard := usecases.UpdateBankCard{
		BankCardRepository: bankCardRepository,
		Crypto:             crypto,
		Log:                log,
	}
	getAllBankCards := usecases.GetAllBankCards{
		BankCardRepository: bankCardRepository,
		Crypto:             crypto,
		Log:                log,
	}

	getAll := usecases.GetAll{
		GetAllTexts:       getAllTexts,
		GetAllBankCards:   getAllBankCards,
		GetAllBinaries:    getAllBinaries,
		GetAllCredentials: getAllCredentials,
		Log:               log,
	}

	return &Application{
		Registration:      registration,
		Login:             login,
		CreateText:        createText,
		UpdateText:        updateText,
		GetAllTexts:       getAllTexts,
		CreateBinary:      createBinary,
		UpdateBinary:      updateBinary,
		GetAllBinaries:    getAllBinaries,
		CreateCredentials: createCredentials,
		UpdateCredentials: updateCredentials,
		GetAllCredentials: getAllCredentials,
		CreateBankCard:    createBankCard,
		UpdateBankCard:    updateBankCard,
		GetAllBankCards:   getAllBankCards,
		GetAll:            getAll,
	}
}
