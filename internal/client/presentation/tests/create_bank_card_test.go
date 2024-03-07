package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestCreateBankCardSuccess(t *testing.T) {
	number := "0000 0000 0000 0000" //nolint: goconst
	validThru := "01/11"            //nolint: goconst
	cvv := "000"
	cardHolder := "name name"

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

	args := []string{
		"gophkeeper",
		"bank_card",
		"create",
		number,
		validThru,
		cvv,
		cardHolder,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	card, err := bankCardRepository.Get(userID, cardID)
	require.NoError(t, err)
	assert.Equal(t, card.Number, number)
	assert.Equal(t, card.ValidThru, validThru)
	assert.Equal(t, card.CVV, cvv)
	assert.Equal(t, card.CardHolder, cardHolder)
}

func TestCreateBankCardNoHolderSuccess(t *testing.T) {
	number := "0000 0000 0000 0000"
	validThru := "01/11"
	cvv := "000"

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

	args := []string{
		"gophkeeper",
		"bank_card",
		"create",
		number,
		validThru,
		cvv,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	card, err := bankCardRepository.Get(userID, cardID)
	require.NoError(t, err)
	assert.Equal(t, card.Number, number)
	assert.Equal(t, card.ValidThru, validThru)
	assert.Equal(t, card.CVV, cvv)
	assert.Equal(t, card.CardHolder, "")
}

func TestCreateBankCardBadRequest(t *testing.T) {
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
		"bank_card",
		"create",
		"0000 0000 0000 0000",
		"01/11",
		"000",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateBankCardUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"bank_card",
		"create",
		"0000 0000 0000 0000",
		"01/11",
		"000",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateBankCardInvalidInput(t *testing.T) {
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

			args := []string{
				"gophkeeper",
				"bank_card",
				"create",
				tt.number,
				tt.validThru,
				tt.cvv,
				tt.cardHolder,
			}

			err = cmd.Run(context.Background(), args)
			require.NoError(t, err)

			_, err = bankCardRepository.Get(userID, cardID)
			require.Error(t, err)
		})
	}
}
