// Package application содержит фабрику приложения и имплементацию CryptoService
package application

import (
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
}

// CreateApplication - Фабрика приложения
func CreateApplication(
	client domain.GophKeeperClientInterface,
	sessionRepository domain.SessionRepositoryInterface,
	textRepository domain.TextRepositoryInterface,
	jwkRepository domain.JWKRepositoryInterface,
) Application {
	jwk, err := jwkRepository.Get()

	if err != nil {
		jwk = nil
	}

	checkToken := usecases.CheckToken{
		Client:        client,
		JWKRepository: jwkRepository,
		Key:           jwk,
	}

	registration := usecases.Registration{
		Client:            client,
		SessionRepository: sessionRepository,
		CheckToken:        &checkToken,
	}

	login := usecases.Login{
		Client:            client,
		SessionRepository: sessionRepository,
		CheckToken:        &checkToken,
	}

	createText := usecases.CreateText{
		Client:         client,
		TextRepository: textRepository,
	}

	updateText := usecases.UpdateText{
		Client:         client,
		TextRepository: textRepository,
	}

	showText := usecases.ShowText{
		CheckToken:     &checkToken,
		TextRepository: textRepository,
	}

	return Application{
		Registration: registration,
		Login:        login,
		CreateText:   createText,
		UpdateText:   updateText,
		ShowText:     showText,
	}
}
