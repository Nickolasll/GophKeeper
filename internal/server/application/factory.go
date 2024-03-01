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
}

// CreateApplication - Фабрика приложения
func CreateApplication(
	userRepository domain.UserRepositoryInterface,
	textRepository domain.TextRepositoryInterface,
	jose services.JOSEService,
	crypto services.CryptoService,
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

	return Application{
		Registration: registration,
		Login:        login,
		CreateText:   createText,
		UpdateText:   updateText,
	}
}
