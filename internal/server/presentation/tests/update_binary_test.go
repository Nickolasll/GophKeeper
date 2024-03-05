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

const binaryURL = "/api/v1/binary/"

func TestUpdateBinaryBadRequest(t *testing.T) { //nolint: dupl
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
			contentType: "multipart/form-data",
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
			token, err := joseService.IssueToken(userID)
			require.NoError(t, err)

			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", binaryURL+tt.resuorceID, bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestUpdateBinarySuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	binID := uuid.New()
	bin := domain.Binary{
		ID:      binID,
		UserID:  userID,
		Content: []byte("my binary message to store"),
	}
	err = binaryRepository.Create(bin)
	require.NoError(t, err)

	message := []byte("my binary message to update")

	bodyReader := bytes.NewReader(message)
	req := httptest.NewRequest("POST", binaryURL+binID.String(), bodyReader)
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	binObj, err := binaryRepository.Get(userID, binID)
	require.NoError(t, err)

	decrypted, err := cryptoService.Decrypt(binObj.Content)
	require.NoError(t, err)

	assert.Equal(t, message, decrypted)
}

func TestUpdateBinaryNotFound(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	message := []byte("my message to update")

	bodyReader := bytes.NewReader(message)
	req := httptest.NewRequest("POST", binaryURL+uuid.NewString(), bodyReader)
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateBinaryInvalidID(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte("message"))
	req := httptest.NewRequest("POST", binaryURL+"invalid", bodyReader)
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
