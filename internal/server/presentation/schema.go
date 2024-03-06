// Package presentation содержит фабрику роутера, обработчики и схемы валидации
package presentation

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validCardNumber *regexp.Regexp
var validThru *regexp.Regexp
var validCVV *regexp.Regexp
var validCardHolder *regexp.Regexp

type registrationPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (registrationPayload) LoadFromJSON(data []byte) (registrationPayload, error) {
	var payload registrationPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}
	err = validate.Struct(payload)

	return payload, err
}

type textPayload struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content" validate:"required,min=1"`
}

func (textPayload) Load(id *uuid.UUID, content []byte) (textPayload, error) {
	payload := textPayload{Content: string(content)}
	if id != nil {
		payload.ID = *id
	}
	err := validate.Struct(payload)

	return payload, err
}

type credentialsPayload struct {
	Name     string `json:"name" validate:"required,min=1"`
	Login    string `json:"login" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

func (credentialsPayload) LoadFromJSON(data []byte) (credentialsPayload, error) {
	var payload credentialsPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}
	err = validate.Struct(payload)

	return payload, err
}

type bankCardPayload struct {
	Number     string `json:"number" validate:"required,is-valid-card-number"`
	ValidThru  string `json:"valid_thru" validate:"required,is-valid-thru"`
	CVV        string `json:"cvv" validate:"required,is-valid-cvv"`
	CardHolder string `json:"card_holder" validate:"is-valid-card-holder"`
}

func (bankCardPayload) LoadFromJSON(data []byte) (bankCardPayload, error) {
	var payload bankCardPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}
	err = validate.Struct(payload)

	return payload, err
}

func validateBankCardNumber(fl validator.FieldLevel) bool {
	return validCardNumber.MatchString(fl.Field().String())
}

func validateValidThru(fl validator.FieldLevel) bool {
	return validThru.MatchString(fl.Field().String())
}

func validateCVV(fl validator.FieldLevel) bool {
	return validCVV.MatchString(fl.Field().String())
}

func validateCardHolder(fl validator.FieldLevel) bool {
	return validCardHolder.MatchString(fl.Field().String())
}

func newValidator() (*validator.Validate, error) {
	validCardNumber = regexp.MustCompile(`\b(\d{4}\s\d{4}\s\d{4}\s\d{4}$)\b`)
	validThru = regexp.MustCompile(`(0[1-9]|1[012])/\d{2}`)
	validCVV = regexp.MustCompile(`^\d{3,4}$`)
	validCardHolder = regexp.MustCompile(`^((?:[A-Za-z]+ ?){0,3})$`)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("is-valid-card-number", validateBankCardNumber)
	if err != nil {
		return validate, err
	}
	err = validate.RegisterValidation("is-valid-thru", validateValidThru)
	if err != nil {
		return validate, err
	}
	err = validate.RegisterValidation("is-valid-cvv", validateCVV)
	if err != nil {
		return validate, err
	}
	err = validate.RegisterValidation("is-valid-card-holder", validateCardHolder)
	if err != nil {
		return validate, err
	}

	return validate, nil
}
