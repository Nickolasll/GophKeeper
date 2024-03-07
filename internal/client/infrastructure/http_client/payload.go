package httpclient

import "encoding/json"

type credentialsPayload struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func credentialsToJSON(name, login, password string) ([]byte, error) {
	cred := credentialsPayload{
		Name:     name,
		Login:    login,
		Password: password,
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
}

func bankCardToJSON(number, validThru, cvv, cardHolder string) ([]byte, error) {
	cred := bankCardPayload{
		Number:     number,
		ValidThru:  validThru,
		CVV:        cvv,
		CardHolder: cardHolder,
	}
	data, err := json.Marshal(cred)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
