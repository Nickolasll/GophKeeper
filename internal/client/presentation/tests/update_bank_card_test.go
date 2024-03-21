package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestUpdateBankCardFullDataSuccess(t *testing.T) {
	number := "0000 0000 0000 0000"
	validThru := "01/11"
	cvv := "000"
	cardHolder := "name name"
	meta := "meta"

	cardID := uuid.New()
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	card := domain.BankCard{
		ID:         cardID,
		Number:     "number",
		ValidThru:  "valid/thru",
		CVV:        "cvv",
		CardHolder: "card holder",
		Meta:       "",
	}

	err = bankCardRepository.Create(userID, &card)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
		"--meta",
		meta,
		"--number",
		number,
		"--valid-thru",
		validThru,
		"--cvv",
		cvv,
		"--card-holder",
		cardHolder,
		cardID.String(),
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	cardObj, err := bankCardRepository.Get(userID, cardID)
	require.NoError(t, err)
	assert.Equal(t, cardObj.Number, number)
	assert.Equal(t, cardObj.ValidThru, validThru)
	assert.Equal(t, cardObj.CVV, cvv)
	assert.Equal(t, cardObj.CardHolder, cardHolder)
	assert.Equal(t, cardObj.Meta, meta)
}

func TestUpdateBankCardInvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		number     string
		validThru  string
		cvv        string
		cardHolder string
	}{
		{
			name:       "bad card number",
			number:     "0000 0000 0000 xxxx",
			validThru:  "01/11",
			cvv:        "000",
			cardHolder: "name name",
		},
		{
			name:       "bad valid thru",
			number:     "0000 0000 0000 0000",
			validThru:  "13/11",
			cvv:        "000",
			cardHolder: "name name",
		},
		{
			name:       "bad cvv",
			number:     "0000 0000 0000 0000",
			validThru:  "01/11",
			cvv:        "xxx",
			cardHolder: "name name",
		},
		{
			name:       "bad card holder name",
			number:     "0000 0000 0000 0000",
			validThru:  "01/11",
			cvv:        "000",
			cardHolder: "name name name name name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardID := uuid.New()
			client := FakeHTTPClient{
				Response: cardID,
			}

			cmd, err := setup(client)
			require.NoError(t, err)
			defer func() {
				err = teardown()
				require.NoError(t, err)
			}()

			userID, err := createSession()
			require.NoError(t, err)

			card := domain.BankCard{
				ID:         cardID,
				Number:     "number",
				ValidThru:  "valid/thru",
				CVV:        "cvv",
				CardHolder: "card holder",
			}
			err = bankCardRepository.Create(userID, &card)
			require.NoError(t, err)

			args := []string{
				"gophkeeper",
				"update",
				"bank-card",
				"--number",
				tt.number,
				"--valid-thru",
				tt.validThru,
				"--cvv",
				tt.cvv,
				"--card-holder",
				tt.cardHolder,
				cardID.String(),
			}

			err = cmd.Run(context.Background(), args)
			require.NoError(t, err)

			cardObj, err := bankCardRepository.Get(userID, cardID)
			require.NoError(t, err)
			assert.NotEqual(t, cardObj.Number, tt.number)
			assert.NotEqual(t, cardObj.ValidThru, tt.validThru)
			assert.NotEqual(t, cardObj.CVV, tt.cvv)
			assert.NotEqual(t, cardObj.CardHolder, tt.cardHolder)
		})
	}
}

func TestUpdateBankCardNoFlags(t *testing.T) {
	client := FakeHTTPClient{
		Err: domain.ErrBadRequest,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	_, err = createSession()
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
		uuid.NewString(),
	}
	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBankCardClientError(t *testing.T) {
	cardID := uuid.New()
	client := FakeHTTPClient{
		Err: domain.ErrBadRequest,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	card := domain.BankCard{
		ID:         cardID,
		Number:     "number",
		ValidThru:  "valid/thru",
		CVV:        "cvv",
		CardHolder: "card holder",
	}
	err = bankCardRepository.Create(userID, &card)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
		"--number",
		"0000 0000 0000 0000",
		cardID.String(),
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBankCardUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBankCardNotFound(t *testing.T) {
	client := FakeHTTPClient{
		Err: domain.ErrBadRequest,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	_, err = createSession()
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
		"--number",
		"0000 0000 0000 0000",
		uuid.NewString(),
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBankCardBadID(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"update",
		"bank-card",
		"--number",
		"0000 0000 0000 0000",
		"invalid",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
