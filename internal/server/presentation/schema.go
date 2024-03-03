// Package presentation содержит фабрику роутера, обработчики и схемы валидации
package presentation

import (
	"encoding/json"

	"github.com/google/uuid"
)

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
