package presentation

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
	if r.Header.Get("Content-Type") != contentType {
		return []byte{}, errInvalidContentType
	}

	return io.ReadAll(r.Body)
}
