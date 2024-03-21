package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
	"github.com/Nickolasll/goph-keeper/internal/server/presentation"
)

const getAllTextsURL = "/api/v1/text/all"

func createText(userID uuid.UUID, message string) (string, error) {
	textID := uuid.New()
	encryptedContent, err := cryptoService.Encrypt([]byte(message))
	if err != nil {
		return "", err
	}
	text := domain.Text{
		ID:      textID,
		UserID:  userID,
		Content: encryptedContent,
	}
	err = textRepository.Create(text)

	return textID.String(), err
}

func TestGetAllTextsSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	firstMessage := "my beautiful text"
	firstID, err := createText(userID, firstMessage)
	require.NoError(t, err)

	secondMessage := "my ugly text"
	secondID, err := createText(userID, secondMessage)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllTextsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllTextsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Data.Texts[0].ID, firstID)
	assert.Equal(t, responseData.Data.Texts[0].Content, firstMessage)

	assert.Equal(t, responseData.Data.Texts[1].ID, secondID)
	assert.Equal(t, responseData.Data.Texts[1].Content, secondMessage)
}

func TestGetAllTextsInternalServerError(t *testing.T) { // nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	text := domain.Text{
		ID:      uuid.New(),
		UserID:  userID,
		Content: []byte("not encrypted"),
	}
	err = textRepository.Create(text)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllTextsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllTextsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Status, false)
	assert.Equal(t, responseData.Message, "cipher: message authentication failed")
}
