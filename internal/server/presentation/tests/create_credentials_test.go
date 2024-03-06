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

func TestCreateCredentialsBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		contentType string
	}{
		{
			name:        "missing field",
			body:        []byte(`{"name": "name", "login": "login"}`),
			contentType: "application/json",
		},
		{
			name:        "wrong fields",
			body:        []byte(`{"field": "value"}`),
			contentType: "application/json",
		},
		{
			name:        "invalid value type",
			body:        []byte(`{"name": "name", "login": "login", "password": 11}`),
			contentType: "application/json",
		},
		{
			name:        "empty string",
			body:        []byte(`{"name":"", "login": "", "password": ""}`),
			contentType: "application/json",
		},
		{
			name:        "not a json",
			body:        []byte(`not a json`),
			contentType: "application/json",
		},
		{
			name:        "wrong content type",
			body:        []byte{},
			contentType: "plain/text",
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
			req := httptest.NewRequest("POST", "/api/v1/credentials/create", bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestCreateCredentialsSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	name := "my cred name"
	login := "login"
	password := "password"

	bodyReader := bytes.NewReader([]byte(`{
		"name": "` + name + `", ` +
		`"login": "` + login + `", ` +
		`"password": "` + password + `"` +
		`}`))
	req := httptest.NewRequest("POST", "/api/v1/credentials/create", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	require.NotEmpty(t, responseRecorder.Header().Get("Location"))
	id := responseRecorder.Header().Get("Location")
	credID, err := uuid.Parse(id)
	require.NoError(t, err)

	cred, err := credentialsRepository.Get(userID, credID)
	require.NoError(t, err)

	decrName, err := cryptoService.Decrypt(cred.Name)
	require.NoError(t, err)

	assert.Equal(t, name, string(decrName))

	decrLogin, err := cryptoService.Decrypt(cred.Login)
	require.NoError(t, err)

	assert.Equal(t, login, string(decrLogin))

	decrPassword, err := cryptoService.Decrypt(cred.Password)
	require.NoError(t, err)

	assert.Equal(t, password, string(decrPassword))
}
