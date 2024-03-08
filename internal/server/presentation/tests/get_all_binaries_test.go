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

const getAllBinariesURL = "/api/v1/binary/all"

func createBinary(userID uuid.UUID, content []byte) (string, error) {
	binID := uuid.New()
	encryptedContent, err := cryptoService.Encrypt(content)
	if err != nil {
		return "", err
	}
	binary := domain.Binary{
		ID:      binID,
		UserID:  userID,
		Content: encryptedContent,
	}
	err = binaryRepository.Create(binary)

	return binID.String(), err
}

func TestGetAllBinariesSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	firstContent := []byte("my beautiful binary")
	firstID, err := createBinary(userID, firstContent)
	require.NoError(t, err)

	secondContent := []byte("my ugly binary")
	secondID, err := createBinary(userID, secondContent)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllBinariesURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllBinariesResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Data.Binaries[0].ID, firstID)
	assert.Equal(t, responseData.Data.Binaries[0].Content, firstContent)

	assert.Equal(t, responseData.Data.Binaries[1].ID, secondID)
	assert.Equal(t, responseData.Data.Binaries[1].Content, secondContent)
}

func TestGetAllBinariesInternalServerError(t *testing.T) { // nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	bin := domain.Binary{
		ID:      uuid.New(),
		UserID:  userID,
		Content: []byte("not encrypted"),
	}
	err = binaryRepository.Create(bin)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllBinariesURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllBinariesResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Status, false)
	assert.Equal(t, responseData.Message, "cipher: message authentication failed")
}
