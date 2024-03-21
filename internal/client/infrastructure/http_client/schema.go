package httpclient

import (
	"encoding/json"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

type credentialsPayload struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

func credentialsToJSON(name, login, password, meta string) ([]byte, error) {
	cred := credentialsPayload{
		Name:     name,
		Login:    login,
		Password: password,
		Meta:     meta,
	}
	data, err := json.Marshal(cred)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

type bankCardPayload struct {
	Number     string `json:"number"`
	ValidThru  string `json:"valid_thru"`
	CVV        string `json:"cvv"`
	CardHolder string `json:"card_holder"`
	Meta       string `json:"meta"`
}

func bankCardToJSON(number, validThru, cvv, cardHolder, meta string) ([]byte, error) {
	cred := bankCardPayload{
		Number:     number,
		ValidThru:  validThru,
		CVV:        cvv,
		CardHolder: cardHolder,
		Meta:       meta,
	}
	data, err := json.Marshal(cred)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

type errorResponse struct {
	Message string `json:"message"`
}

type getAllTextsResponse struct {
	Data struct {
		Texts []domain.Text `json:"texts"`
	} `json:"data"`
}

type getAllBinariesResponse struct {
	Data struct {
		Binaries []domain.Binary `json:"binaries"`
	} `json:"data"`
}

type getAllCredentialsResponse struct {
	Data struct {
		Credentials []domain.Credentials `json:"credentials"`
	} `json:"data"`
}

type getAllBankCardsResponse struct {
	Data struct {
		BankCards []domain.BankCard `json:"bank_cards"`
	} `json:"data"`
}

type getAllResponse struct {
	Data struct {
		Texts       []domain.Text        `json:"texts"`
		Binaries    []domain.Binary      `json:"binaries"`
		Credentials []domain.Credentials `json:"credentials"`
		BankCards   []domain.BankCard    `json:"bank_cards"`
	} `json:"data"`
}
