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

func TestLoginNoUserUnauthorized(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	bodyReader := bytes.NewReader([]byte(`{"login": "` + uuid.NewString() + `", "password": "password"}`)) //nolint: goconst
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestLoginSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	login := uuid.NewString()
	password := "password"
	hashedPassword := jose.Hash(password)

	user := domain.User{
		ID:       uuid.New(),
		Login:    login,
		Password: hashedPassword,
	}
	err = userRepository.Create(user)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte(`{"login": "` + login + `", "password": "` + password + `"}`))
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Header().Get("Authorization"))
}

func TestLoginWrongPasswordUnauthorized(t *testing.T) { //nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	login := uuid.NewString()
	user := domain.User{
		ID:       uuid.New(),
		Login:    login,
		Password: "password",
	}
	err = userRepository.Create(user)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte(`{"login": "` + login + `", "password": "qwerty"}`))
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}
