package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTextBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		contentType string
	}{
		{
			name:        "no content",
			body:        []byte{},
			contentType: "plain/text",
		},
		{
			name:        "wrong content type",
			body:        []byte{},
			contentType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := setup()
			require.NoError(t, err)
			defer teardown()

			userID := uuid.New()
			err = createUser(userID)
			require.NoError(t, err)
			token, err := joseService.IssueToken(userID)
			require.NoError(t, err)

			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/text/create", bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestCreateTextSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	message := "my text message to store"

	bodyReader := bytes.NewReader([]byte(message))
	req := httptest.NewRequest("POST", "/api/v1/text/create", bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	require.NotEmpty(t, responseRecorder.Header().Get("Location"))
	id := responseRecorder.Header().Get("Location")
	textID, err := uuid.Parse(id)
	require.NoError(t, err)

	textObj, err := textRepository.Get(userID, textID)
	require.NoError(t, err)

	decrypted, err := cryptoService.Decrypt(textObj.Content)
	require.NoError(t, err)

	assert.Equal(t, message, string(decrypted))
}
