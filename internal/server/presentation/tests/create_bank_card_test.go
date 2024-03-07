// nolint: goconst
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

func TestCreateBankCardBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		contentType string
	}{
		{
			name:        "missing field",
			body:        []byte(`{"number": "1234 5678 1234 5678", "cvv": "123"}`),
			contentType: "application/json",
		},
		{
			name:        "wrong fields",
			body:        []byte(`{"field": "value"}`),
			contentType: "application/json",
		},
		{
			name:        "invalid value type",
			body:        []byte(`{"number": 11, "cvv": 11, "valid_thru": 11}`),
			contentType: "application/json",
		},
		{
			name:        "empty string",
			body:        []byte(`{"number":"", "cvv": "", "valid_thru": ""}`),
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
		{
			name:        "invalid card number",
			body:        []byte(`{"number": "xxxx 5678 1234 5678", "cvv": "123", "valid_thru": "01/30"}`),
			contentType: "application/json",
		},
		{
			name:        "invalid card cvv",
			body:        []byte(`{"number": "1234 5678 1234 5678", "cvv": "12356", "valid_thru": "01/30"}`),
			contentType: "application/json",
		},
		{
			name:        "invalid card valid_thru",
			body:        []byte(`{"number": "1234 5678 1234 5678", "cvv": "123", "valid_thru": "15/30"}`),
			contentType: "application/json",
		},
		{
			name:        "invalid card holder",
			body:        []byte(`{"number": "1234 5678 1234 5678", "cvv": "123", "valid_thru": "01/30", "card_holder": "name name name name"}`),
			contentType: "application/json",
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
			req := httptest.NewRequest("POST", "/api/v1/bank_card/create", bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestCreateBankCardSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	number := "1234 5678 1234 5678"
	validThru := "01/30"
	cvv := "123"
	cardHolder := "Card Holder"

	bodyReader := bytes.NewReader([]byte(`{
		"number": "` + number + `", ` +
		`"valid_thru": "` + validThru + `", ` +
		`"cvv": "` + cvv + `", ` +
		`"card_holder": "` + cardHolder + `"` +
		`}`))
	req := httptest.NewRequest("POST", "/api/v1/bank_card/create", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	require.NotEmpty(t, responseRecorder.Header().Get("Location"))
	id := responseRecorder.Header().Get("Location")
	cardID, err := uuid.Parse(id)
	require.NoError(t, err)

	card, err := cardRepository.Get(userID, cardID)
	require.NoError(t, err)

	decrNumber, err := cryptoService.Decrypt(card.Number)
	require.NoError(t, err)

	assert.Equal(t, number, string(decrNumber))

	decrValidThru, err := cryptoService.Decrypt(card.ValidThru)
	require.NoError(t, err)

	assert.Equal(t, validThru, string(decrValidThru))

	decrCVV, err := cryptoService.Decrypt(card.CVV)
	require.NoError(t, err)

	assert.Equal(t, cvv, string(decrCVV))

	decrCardHolder, err := cryptoService.Decrypt(card.CardHolder)
	require.NoError(t, err)

	assert.Equal(t, cardHolder, string(decrCardHolder))
}

func TestCreateBankCardNoHolderSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	number := "1234 5678 1234 5678"
	validThru := "01/30"
	cvv := "123"

	bodyReader := bytes.NewReader([]byte(`{
		"number": "` + number + `", ` +
		`"valid_thru": "` + validThru + `", ` +
		`"cvv": "` + cvv + `"` +
		`}`))
	req := httptest.NewRequest("POST", "/api/v1/bank_card/create", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	require.NotEmpty(t, responseRecorder.Header().Get("Location"))
	id := responseRecorder.Header().Get("Location")
	credID, err := uuid.Parse(id)
	require.NoError(t, err)

	card, err := cardRepository.Get(userID, credID)
	require.NoError(t, err)

	decrNumber, err := cryptoService.Decrypt(card.Number)
	require.NoError(t, err)

	assert.Equal(t, number, string(decrNumber))

	decrValidThru, err := cryptoService.Decrypt(card.ValidThru)
	require.NoError(t, err)

	assert.Equal(t, validThru, string(decrValidThru))

	decrCVV, err := cryptoService.Decrypt(card.CVV)
	require.NoError(t, err)

	assert.Equal(t, cvv, string(decrCVV))

	decrCardHolder, err := cryptoService.Decrypt(card.CardHolder)
	require.NoError(t, err)

	assert.Equal(t, "", string(decrCardHolder))
}
