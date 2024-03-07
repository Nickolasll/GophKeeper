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

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

const cardURL = "/api/v1/bank_card/"

func TestUpdateBankCardBadRequest(t *testing.T) { //nolint: dupl
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
			req := httptest.NewRequest("POST", cardURL+tt.resuorceID, bodyReader)
			req.Header.Add("Content-Type", tt.contentType)
			req.Header.Add("Authorization", string(token))
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, req)
			assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		})
	}
}

func TestUpdateBankCardSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	cardID := uuid.New()
	card := domain.BankCard{
		ID:         cardID,
		UserID:     userID,
		Number:     []byte("0000 0000 0000 0000"),
		CVV:        []byte("000"),
		ValidThru:  []byte("01/11"),
		CardHolder: []byte(""),
	}
	err = cardRepository.Create(&card)
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
	req := httptest.NewRequest("POST", cardURL+cardID.String(), bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	cardObj, err := cardRepository.Get(userID, cardID)
	require.NoError(t, err)

	decrNumber, err := cryptoService.Decrypt(cardObj.Number)
	require.NoError(t, err)

	assert.Equal(t, number, string(decrNumber))

	decrValidThru, err := cryptoService.Decrypt(cardObj.ValidThru)
	require.NoError(t, err)

	assert.Equal(t, validThru, string(decrValidThru))

	decrCVV, err := cryptoService.Decrypt(cardObj.CVV)
	require.NoError(t, err)

	assert.Equal(t, cvv, string(decrCVV))

	decrCardHolder, err := cryptoService.Decrypt(cardObj.CardHolder)
	require.NoError(t, err)

	assert.Equal(t, cardHolder, string(decrCardHolder))
}

func TestUpdateBankCardNotFound(t *testing.T) {
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
	req := httptest.NewRequest("POST", cardURL+uuid.NewString(), bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateBankCardInvalidID(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	bodyReader := bytes.NewReader([]byte{})
	req := httptest.NewRequest("POST", cardURL+"invalid", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
