package presentation

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validCardNumber, validThru, validCVV, validCardHolder *regexp.Regexp

type registrationPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (registrationPayload) Load(data []byte) (registrationPayload, error) {
	var payload registrationPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}
	err = validate.Struct(payload)

	return payload, err
}

type credentialsPayload struct {
	Name     string `json:"name" validate:"required,min=1"`
	Login    string `json:"login" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
	Meta     string `json:"meta"`
}

func (credentialsPayload) Load(data []byte) (credentialsPayload, error) {
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
	Meta       string `json:"meta"`
}

func (bankCardPayload) Load(data []byte) (bankCardPayload, error) {
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

type textResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type GetAllTextsResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Texts []textResponse `json:"texts"`
	} `json:"data"`
}

type binaryResponse struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`
}

type GetAllBinariesResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Binaries []binaryResponse `json:"binaries"`
	} `json:"data"`
}

type credentialsResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

type GetAllCredentialsResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Credentials []credentialsResponse `json:"credentials"`
	} `json:"data"`
}

type bankCardResponse struct {
	ID         string `json:"id"`
	Number     string `json:"number"`
	ValidThru  string `json:"valid_thru"`
	CVV        string `json:"cvv"`
	CardHolder string `json:"card_holder"`
	Meta       string `json:"meta"`
}

type GetAllBankCardsResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		BankCards []bankCardResponse `json:"bank_cards"`
	} `json:"data"`
}

type GetAllResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Texts       []textResponse        `json:"texts"`
		Binaries    []binaryResponse      `json:"binaries"`
		Credentials []credentialsResponse `json:"credentials"`
		BankCards   []bankCardResponse    `json:"bank_cards"`
	} `json:"data"`
}

type ErrorResponse struct {
	Status  bool     `json:"status"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
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
