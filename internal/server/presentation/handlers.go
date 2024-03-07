// Package presentation содержит фабрику роутера, обработчики и схемы валидации
package presentation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

func registrationHandler(w http.ResponseWriter, r *http.Request) { //nolint: dupl
	var payload registrationPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Registration.Do(payload.Login, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrLoginAlreadyInUse) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}

	w.Header().Set("Authorization", string(token))
	w.WriteHeader(http.StatusOK)
}

func loginHandler(w http.ResponseWriter, r *http.Request) { //nolint: dupl
	var payload registrationPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Login.Do(payload.Login, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrLoginOrPasswordIsInvalid) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}

	w.Header().Set("Authorization", string(token))
	w.WriteHeader(http.StatusOK)
}

func createTextHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	body, err := parseBody(textType, r)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	textID, err := app.CreateText.Do(userID, string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", textID.String())
	w.WriteHeader(http.StatusCreated)
}

func updateTextHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	id, err := getRouteID(r, "textID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}

	body, err := parseBody(textType, r)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateText.Do(userID, id, string(body))
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

func getCertsHandler(w http.ResponseWriter, _ *http.Request) {
	certs, err := joseService.GetCerts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	resp, err := json.Marshal(certs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}
}

// Health godoc
// @Summary Запрос состояния сервиса
// @ID health
// @Success 200
// @Failure 500
// @Router /health [get]
func getHealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func createBinaryHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	body, err := parseBody(binaryType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	binID, err := app.CreateBinary.Do(userID, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", binID.String())
	w.WriteHeader(http.StatusCreated)
}

func updateBinaryHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	id, err := getRouteID(r, "binaryID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(binaryType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateBinary.Do(userID, id, body)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

func createCredentialsHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload credentialsPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	credID, err := app.CreateCredentials.Do(
		userID,
		payload.Name,
		payload.Login,
		payload.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", credID.String())
	w.WriteHeader(http.StatusCreated)
}

func updateCredentialsHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload credentialsPayload
	id, err := getRouteID(r, "credID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateCredentials.Do(
		userID,
		id,
		payload.Name,
		payload.Login,
		payload.Password,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

func createBankCardHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload bankCardPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	cardID, err := app.CreateBankCard.Do(
		userID,
		payload.Number,
		payload.ValidThru,
		payload.CVV,
		payload.CardHolder,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", cardID.String())
	w.WriteHeader(http.StatusCreated)
}

func updateBankCardHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload bankCardPayload
	id, err := getRouteID(r, "cardID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateBankCard.Do(
		userID,
		id,
		payload.Number,
		payload.ValidThru,
		payload.CVV,
		payload.CardHolder,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}
