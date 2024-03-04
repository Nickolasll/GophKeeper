// Package application содержит фабрику приложения
package application

import (
	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
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
	// CreateText - Сценарий использования для создания зашифрованных бинарных данных
	CreateBinary usecases.CreateBinary
	// UpdateText - Сценарий использования для обновления существующих зашифрованных бинарных данных
	UpdateBinary usecases.UpdateBinary
}

// New - Фабрика приложения
func New(
	jose services.JOSEService,
	crypto domain.CryptoServiceInterface,
	userRepository domain.UserRepositoryInterface,
	textRepository domain.TextRepositoryInterface,
	binaryRepository domain.BinaryRepositoryInterface,
) Application {
	registration := usecases.Registration{
		UserRepository: userRepository,
		JOSE:           jose,
	}

	login := usecases.Login{
		UserRepository: userRepository,
		JOSE:           jose,
	}

	createText := usecases.CreateText{
		TextRepository: textRepository,
		Crypto:         crypto,
	}

	updateText := usecases.UpdateText{
		TextRepository: textRepository,
		Crypto:         crypto,
	}

	createBinary := usecases.CreateBinary{
		BinaryRepository: binaryRepository,
		Crypto:           crypto,
	}

	updateBinary := usecases.UpdateBinary{
		BinaryRepository: binaryRepository,
		Crypto:           crypto,
	}

	return Application{
		Registration: registration,
		Login:        login,
		CreateText:   createText,
		UpdateText:   updateText,
		CreateBinary: createBinary,
		UpdateBinary: updateBinary,
	}
}
