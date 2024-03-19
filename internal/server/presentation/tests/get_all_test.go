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

const getAllURL = "/api/v1/all"

func TestGetAllSuccess(t *testing.T) { // nolint: funlen
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	textMessage := "my beautiful text"
	firstTextID, err := createText(userID, textMessage)
	require.NoError(t, err)

	secondTextID, err := createText(userID, textMessage)
	require.NoError(t, err)

	name := "my credentials name"
	login := "login"
	password := "password"
	credMeta := "my credentials meta"

	firstCredID, err := createCredentials(userID, name, login, password, credMeta)
	require.NoError(t, err)

	secondCredID, err := createCredentials(userID, name, login, password, credMeta)
	require.NoError(t, err)

	binaryContent := []byte("my beautiful binary")
	firstBinaryID, err := createBinary(userID, binaryContent)
	require.NoError(t, err)

	secondBinaryID, err := createBinary(userID, binaryContent)
	require.NoError(t, err)

	number := "0000 0000 0000 0000"
	validThru := "01/11"
	cvv := "000"
	cardHolder := "name name"
	cardMeta := "my card meta"

	firstCardID, err := createBankCard(userID, number, validThru, cvv, cardHolder, cardMeta)
	require.NoError(t, err)

	secondCardID, err := createBankCard(userID, number, validThru, cvv, cardHolder, cardMeta)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Data.Texts[0].ID, firstTextID)
	assert.Equal(t, responseData.Data.Texts[0].Content, textMessage)

	assert.Equal(t, responseData.Data.Texts[1].ID, secondTextID)
	assert.Equal(t, responseData.Data.Texts[1].Content, textMessage)

	assert.Equal(t, responseData.Data.Credentials[0].ID, firstCredID)
	assert.Equal(t, responseData.Data.Credentials[0].Name, name)
	assert.Equal(t, responseData.Data.Credentials[0].Login, login)
	assert.Equal(t, responseData.Data.Credentials[0].Password, password)
	assert.Equal(t, responseData.Data.Credentials[0].Meta, credMeta)

	assert.Equal(t, responseData.Data.Credentials[1].ID, secondCredID)
	assert.Equal(t, responseData.Data.Credentials[1].Name, name)
	assert.Equal(t, responseData.Data.Credentials[1].Login, login)
	assert.Equal(t, responseData.Data.Credentials[1].Password, password)
	assert.Equal(t, responseData.Data.Credentials[1].Meta, credMeta)

	assert.Equal(t, responseData.Data.Binaries[0].ID, firstBinaryID)
	assert.Equal(t, responseData.Data.Binaries[0].Content, binaryContent)

	assert.Equal(t, responseData.Data.Binaries[1].ID, secondBinaryID)
	assert.Equal(t, responseData.Data.Binaries[1].Content, binaryContent)

	assert.Equal(t, responseData.Data.BankCards[0].ID, firstCardID)
	assert.Equal(t, responseData.Data.BankCards[0].Number, number)
	assert.Equal(t, responseData.Data.BankCards[0].ValidThru, validThru)
	assert.Equal(t, responseData.Data.BankCards[0].CVV, cvv)
	assert.Equal(t, responseData.Data.BankCards[0].CardHolder, cardHolder)
	assert.Equal(t, responseData.Data.BankCards[0].Meta, cardMeta)

	assert.Equal(t, responseData.Data.BankCards[1].ID, secondCardID)
	assert.Equal(t, responseData.Data.BankCards[1].Number, number)
	assert.Equal(t, responseData.Data.BankCards[1].ValidThru, validThru)
	assert.Equal(t, responseData.Data.BankCards[1].CVV, cvv)
	assert.Equal(t, responseData.Data.BankCards[1].Meta, cardMeta)
}

func TestGetAllInternalServerError(t *testing.T) { // nolint: dupl
	router, err := setup()
	require.NoError(t, err)
	defer teardown()

	userID := uuid.New()
	err = createUser(userID)
	require.NoError(t, err)
	token, err := joseService.IssueToken(userID)
	require.NoError(t, err)

	text := domain.Text{
		ID:      uuid.New(),
		UserID:  userID,
		Content: []byte("not encrypted"),
	}
	err = textRepository.Create(text)
	require.NoError(t, err)

	bodyReader := bytes.NewReader(nil)
	req := httptest.NewRequest("GET", getAllURL, bodyReader)
	req.Header.Add("Authorization", string(token))
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	responseData := presentation.GetAllResponse{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseData)
	require.NoError(t, err)

	assert.Equal(t, responseData.Status, false)
	assert.Equal(t, responseData.Message, "cipher: message authentication failed")
}
