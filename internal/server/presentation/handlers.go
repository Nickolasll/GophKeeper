// Package presentation содержит фабрику роутера, обработчики и схемы валидации
package presentation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

const jsonType = "application/json"
const textType = "plain/text"

type authenticatedHandler func(w http.ResponseWriter, r *http.Request, userID uuid.UUID)

var errInvalidContentType = errors.New("invalid content type")

func getRouteID(r *http.Request, name string) (uuid.UUID, error) {
	strID := chi.URLParam(r, name)
	id, err := uuid.Parse(strID)

	return id, err
}

func parseBody(contentType string, r *http.Request) ([]byte, error) {
	if r.Header.Get("Content-Type") != contentType {
		return []byte{}, errInvalidContentType
	}

	return io.ReadAll(r.Body)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) { //nolint: dupl
	var payload registrationPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.LoadFromJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Registration.Execute(payload.Login, payload.Password)
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
	payload, err = payload.LoadFromJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Login.Execute(payload.Login, payload.Password)
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
	var payload textPayload
	body, err := parseBody(textType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(nil, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	textID, err := app.CreateText.Execute(userID, payload.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", textID.String())
	w.WriteHeader(http.StatusCreated)
}

func updateTextHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload textPayload
	id, err := getRouteID(r, "textID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}

	body, err := parseBody(textType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(&id, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateText.Execute(userID, payload.ID, payload.Content)
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
	certs, err := jose.GetCerts()
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
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}
}
