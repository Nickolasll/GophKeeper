package presentation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const contentTypeHeader = "Content-Type"
const jsonType = "application/json"
const textType = "plain/text"
const binaryType = "multipart/form-data"

var errInvalidContentType = errors.New("invalid content type")

type authenticatedHandler func(w http.ResponseWriter, r *http.Request, userID uuid.UUID)

func getRouteID(r *http.Request, name string) (uuid.UUID, error) {
	strID := chi.URLParam(r, name)
	id, err := uuid.Parse(strID)

	return id, err
}

func parseBody(contentType string, r *http.Request) ([]byte, error) {
	if r.Header.Get(contentTypeHeader) != contentType {
		return []byte{}, errInvalidContentType
	}

	return io.ReadAll(r.Body)
}

func makeResponse(w http.ResponseWriter, statusCode int, response any) error {
	w.WriteHeader(statusCode)
	responseData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(responseData)
	if err != nil {
		return err
	}

	return nil
}

func responseError(w http.ResponseWriter, message string) error {
	errResp := ErrorResponse{
		Status:  false,
		Message: message,
	}
	err := makeResponse(w, http.StatusInternalServerError, errResp)
	if err != nil {
		return err
	}

	return nil
}
