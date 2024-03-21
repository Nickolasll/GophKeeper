package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHealth(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	req := httptest.NewRequest("GET", "/api/v1/health", http.NoBody)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}
