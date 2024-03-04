package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

const textURL = "/api/v1/text/"

func TestUpdateTextBadRequest(t *testing.T) { //nolint: dupl
	tests := []struct {
		name        string
		body        []byte
		contentType string
		resuorceID  string
	}{
		{
			name:        "no content",
			body:        []byte{},
			contentType: "plain/text",
			resuorceID:  uuid.NewString(),
		},
		{
			name:        "wrong content type",
			body:        []byte{},
			contentType: "application/json",
			resuorceID:  uuid.NewString(),
		},
		{
			name:        "wrong resource type",
			body:        []byte{},
			contentType: "application/json",
			resuorceID:  "not_a_UUID",
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
			token, err := jose.IssueToken(userID)
			require.NoError(t, err)

			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", textURL+tt.resuorceID, bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestUpdateTextSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := jose.IssueToken(userID)
	require.NoError(t, err)

	textID := uuid.New()
	text := domain.Text{
		ID:      textID,
		UserID:  userID,
		Content: []byte("my text message to store"),
	}
	err = textRepository.Create(text)
	require.NoError(t, err)

	message := "my message to update"

	bodyReader := bytes.NewReader([]byte(message))
	req := httptest.NewRequest("POST", textURL+textID.String(), bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	textObj, err := textRepository.Get(userID, textID)
	require.NoError(t, err)

	decrypted, err := cryptoService.Decrypt(textObj.Content)
	require.NoError(t, err)

	assert.Equal(t, message, string(decrypted))
}

func TestUpdateTextNotFound(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := jose.IssueToken(userID)
	require.NoError(t, err)

	message := "my message to update"

	bodyReader := bytes.NewReader([]byte(message))
	req := httptest.NewRequest("POST", textURL+uuid.NewString(), bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateTextInvalidID(t *testing.T) { //nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := jose.IssueToken(userID)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte("message"))
	req := httptest.NewRequest("POST", textURL+"invalid", bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
