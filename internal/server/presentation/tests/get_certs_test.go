package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCerts(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	req := httptest.NewRequest("GET", "/api/v1/auth/certs", http.NoBody)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, body)

	userID := uuid.New()
	token, err := jose.IssueToken(userID)
	require.NoError(t, err)

	key, err := jwk.ParseKey(body)
	require.NoError(t, err)

	tok, err := jwt.Parse(
		token,
		jwt.WithKey(jwa.HS256, key),
		jwt.WithValidate(true),
	)

	require.NoError(t, err)
	v, _ := tok.Get("UserID")
	str, _ := v.(string)
	userIDfromToken, err := uuid.Parse(str)
	require.NoError(t, err)
	assert.Equal(t, userID, userIDfromToken)

	serializedKey, err := json.Marshal(key)
	require.NoError(t, err)
	parsedKey, err := jwk.ParseKey(serializedKey)
	require.NoError(t, err)
	parsedToken, err := jwt.Parse(
		token,
		jwt.WithKey(jwa.HS256, parsedKey),
		jwt.WithValidate(true),
	)
	require.NoError(t, err)
	v, _ = parsedToken.Get("UserID")
	str, _ = v.(string)
	userIDfromToken, err = uuid.Parse(str)
	require.NoError(t, err)
	assert.Equal(t, userID, userIDfromToken)
}
