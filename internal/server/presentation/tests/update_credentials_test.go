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

const credURL = "/api/v1/credentials/"

func TestUpdateCredentialsBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		contentType string
		resuorceID  string
	}{
		{
			name:        "no content",
			body:        []byte{},
			contentType: "application/json",
			resuorceID:  uuid.NewString(),
		},
		{
			name:        "wrong content type",
			body:        []byte{},
			contentType: "plain/text",
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
			token, err := joseService.IssueToken(userID)
			require.NoError(t, err)

			bodyReader := bytes.NewReader(tt.body)
			req := httptest.NewRequest("POST", credURL+tt.resuorceID, bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestUpdateCredentialsSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	credID := uuid.New()
	cred := domain.Credentials{
		ID:       credID,
		UserID:   userID,
		Name:     []byte("name"),
		Login:    []byte("login"),
		Password: []byte("password"),
	}
	err = credentialsRepository.Create(&cred)
	require.NoError(t, err)

	name := "my name to update"
	login := "my login to update"
	password := "my password to update"

	bodyReader := bytes.NewReader([]byte(`{
		"name": "` + name + `", ` +
		`"login": "` + login + `", ` +
		`"password": "` + password + `"` +
		`}`))
	req := httptest.NewRequest("POST", credURL+credID.String(), bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	credObj, err := credentialsRepository.Get(userID, credID)
	require.NoError(t, err)

	decrName, err := cryptoService.Decrypt(credObj.Name)
	require.NoError(t, err)

	assert.Equal(t, name, string(decrName))

	decrLogin, err := cryptoService.Decrypt(credObj.Login)
	require.NoError(t, err)

	assert.Equal(t, login, string(decrLogin))

	decrPassword, err := cryptoService.Decrypt(credObj.Password)
	require.NoError(t, err)

	assert.Equal(t, password, string(decrPassword))
}

func TestUpdateCredentialsNotFound(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	name := "my name to update"
	login := "my login to update"
	password := "my password to update"

	bodyReader := bytes.NewReader([]byte(`{
		"name": "` + name + `", ` +
		`"login": "` + login + `", ` +
		`"password": "` + password + `"` +
		`}`))
	req := httptest.NewRequest("POST", credURL+uuid.NewString(), bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateCredentialsInvalidID(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	name := "my name to update"
	login := "my login to update"
	password := "my password to update"

	bodyReader := bytes.NewReader([]byte(`{
		"name": "` + name + `", ` +
		`"login": "` + login + `", ` +
		`"password": "` + password + `"` +
		`}`))
	req := httptest.NewRequest("POST", textURL+"invalid", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
