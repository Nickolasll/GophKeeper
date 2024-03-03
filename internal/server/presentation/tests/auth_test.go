package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationNoTokenValue(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	req := httptest.NewRequest("POST", "/api/v1/text/create", http.NoBody)
	req.Header.Add("Authorization", "")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthorizationInvalidTokenValue(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	req := httptest.NewRequest("POST", "/api/v1/text/create", http.NoBody)
	req.Header.Add("Authorization", "invalid token value")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthorizationExpiredToken(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	jose.TokenExp = 0
	expiredToken, err := jose.IssueToken(userID)
	require.NoError(t, err)
	jose.TokenExp = 30

	req := httptest.NewRequest("POST", "/api/v1/text/create", http.NoBody)
	req.Header.Add("Authorization", string(expiredToken))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}
