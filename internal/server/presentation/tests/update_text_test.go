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
	req := httptest.NewRequest("POST", "/api/v1/text/"+textID.String(), bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	textObj, err := textRepository.Get(textID, userID)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(textObj.Content)
	require.NoError(t, err)

	assert.Equal(t, message, decrypted)
}

func TestUpdateNotFound(t *testing.T) {
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
	req := httptest.NewRequest("POST", "/api/v1/text/"+uuid.NewString(), bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateInvalidID(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := jose.IssueToken(userID)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte("message"))
	req := httptest.NewRequest("POST", "/api/v1/text/invalid", bodyReader)
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
