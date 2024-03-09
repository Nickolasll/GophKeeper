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

const getAllBankCardsURL = "/api/v1/bank_card/all"

func createBankCard(userID uuid.UUID, number, validThru, cvv, cardHolder string) (string, error) { // nolint: unparam
	cardID := uuid.New()
	encryptedNumber, err := cryptoService.Encrypt([]byte(number))
	if err != nil {
		return "", err
	}
	encryptedValidThru, err := cryptoService.Encrypt([]byte(validThru))
	if err != nil {
		return "", err
	}
	encryptedCVV, err := cryptoService.Encrypt([]byte(cvv))
	if err != nil {
		return "", err
	}
	encryptedCardHolder, err := cryptoService.Encrypt([]byte(cardHolder))
	if err != nil {
		return "", err
	}
	card := domain.BankCard{
		ID:         cardID,
		UserID:     userID,
		Number:     encryptedNumber,
		ValidThru:  encryptedValidThru,
		CardHolder: encryptedCardHolder,
		CVV:        encryptedCVV,
	}
	err = cardRepository.Create(&card)

	return cardID.String(), err
}

func TestGetAllBankCardsSuccess(t *testing.T) {
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	number := "0000 0000 0000 0000"
	validThru := "01/11"
	cvv := "000"
	cardHolder := "name name"

	firstID, err := createBankCard(userID, number, validThru, cvv, cardHolder)
	require.NoError(t, err)

	secondID, err := createBankCard(userID, number, validThru, cvv, cardHolder)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllBankCardsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllBankCardsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Data.BankCards[0].ID, firstID)
	assert.Equal(t, responseData.Data.BankCards[0].Number, number)
	assert.Equal(t, responseData.Data.BankCards[0].ValidThru, validThru)
	assert.Equal(t, responseData.Data.BankCards[0].CVV, cvv)
	assert.Equal(t, responseData.Data.BankCards[0].CardHolder, cardHolder)

	assert.Equal(t, responseData.Data.BankCards[1].ID, secondID)
	assert.Equal(t, responseData.Data.BankCards[1].Number, number)
	assert.Equal(t, responseData.Data.BankCards[1].ValidThru, validThru)
	assert.Equal(t, responseData.Data.BankCards[1].CVV, cvv)
	assert.Equal(t, responseData.Data.BankCards[1].CardHolder, cardHolder)
}

func TestGetAllBankCardsInternalServerError(t *testing.T) { // nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	notEncrypted := []byte("not encrypted")
	card := domain.BankCard{
		ID:         uuid.New(),
		UserID:     userID,
		Number:     notEncrypted,
		ValidThru:  notEncrypted,
		CardHolder: notEncrypted,
		CVV:        notEncrypted,
	}

	err = cardRepository.Create(&card)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllBankCardsURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllBankCardsResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Status, false)
	assert.Equal(t, responseData.Message, "cipher: message authentication failed")
}
