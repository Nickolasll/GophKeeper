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

const getAllCredentialsURL = "/api/v1/credentials/all" // nolint: gosec

func createCredentials(userID uuid.UUID, name, login, password string) (string, error) {
	credID := uuid.New()
	encryptedName, err := cryptoService.Encrypt([]byte(name))
	if err != nil {
		return "", err
	}
	encryptedLogin, err := cryptoService.Encrypt([]byte(login))
	if err != nil {
		return "", err
	}
	encryptedPassword, err := cryptoService.Encrypt([]byte(password))
	if err != nil {
		return "", err
	}
	cred := domain.Credentials{
		ID:       credID,
		UserID:   userID,
		Name:     encryptedName,
		Login:    encryptedLogin,
		Password: encryptedPassword,
	}
	err = credentialsRepository.Create(&cred)

	return credID.String(), err
}

func TestGetAllCredentialsSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	name := "my credentials name"
	login := "login"
	password := "password"

	firstID, err := createCredentials(userID, name, login, password)
	require.NoError(t, err)

	secondID, err := createCredentials(userID, name, login, password)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllCredentialsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllCredentialsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Data.Credentials[0].ID, firstID)
	assert.Equal(t, responseData.Data.Credentials[0].Name, name)
	assert.Equal(t, responseData.Data.Credentials[0].Login, login)
	assert.Equal(t, responseData.Data.Credentials[0].Password, password)

	assert.Equal(t, responseData.Data.Credentials[1].ID, secondID)
	assert.Equal(t, responseData.Data.Credentials[1].Name, name)
	assert.Equal(t, responseData.Data.Credentials[1].Login, login)
	assert.Equal(t, responseData.Data.Credentials[1].Password, password)
}

func TestGetAllCredentialsInternalServerError(t *testing.T) { // nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	cred := domain.Credentials{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     []byte("not encrypted"),
		Login:    []byte("not encrypted"),
		Password: []byte("not encrypted"),
	}
	err = credentialsRepository.Create(&cred)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllCredentialsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllCredentialsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Status, false)
	assert.Equal(t, responseData.Message, "cipher: message authentication failed")
}
