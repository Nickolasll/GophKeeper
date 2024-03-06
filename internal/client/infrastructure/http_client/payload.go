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
