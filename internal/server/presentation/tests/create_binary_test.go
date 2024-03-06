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

func TestBinaryBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		contentType string
	}{
		{
			name:        "no content",
			body:        []byte{},
			contentType: "multipart/form-data",
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
			req := httptest.NewRequest("POST", "/api/v1/binary/create", bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestBinarySuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	message := []byte("my secret binary message")

	bodyReader := bytes.NewReader(message)
	req := httptest.NewRequest("POST", "/api/v1/binary/create", bodyReader)
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	require.NotEmpty(t, responseRecorder.Header().Get("Location"))
	id := responseRecorder.Header().Get("Location")
	binID, err := uuid.Parse(id)
	require.NoError(t, err)

	binObj, err := binaryRepository.Get(userID, binID)
	require.NoError(t, err)

	decrypted, err := cryptoService.Decrypt(binObj.Content)
	require.NoError(t, err)

	assert.Equal(t, message, decrypted)
}
