package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestSyncBankCardsOverrideSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.BankCard{
			{
				ID:         uuid.New(),
				Number:     "0000 0000 0000 0000",
				ValidThru:  "01/11",
				CVV:        "000",
				CardHolder: "name name",
			},
			{
				ID:         uuid.New(),
				Number:     "0000 0000 0000 0000",
				ValidThru:  "01/11",
				CVV:        "000",
				CardHolder: "name name",
			},
		},
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
		ID:         uuid.New(),
		Number:     "1111 2222 3333 4444",
		ValidThru:  "11/11",
		CVV:        "1111",
		CardHolder: "not name",
	}
	err = bankCardRepository.Create(userID, &card)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"bank_card",
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	cards, err := bankCardRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(cards), 2)
	for _, v := range cards {
		assert.NotEqual(t, v.ID, card.ID)
		assert.NotEqual(t, v.Number, card.Number)
		assert.NotEqual(t, v.ValidThru, card.ValidThru)
		assert.NotEqual(t, v.CVV, card.CVV)
		assert.NotEqual(t, v.CardHolder, card.CardHolder)
	}
}

func TestSyncBankCardsSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.BankCard{
			{
				ID:         uuid.New(),
				Number:     "0000 0000 0000 0000",
				ValidThru:  "01/11",
				CVV:        "000",
				CardHolder: "name name",
			},
			{
				ID:         uuid.New(),
				Number:     "0000 0000 0000 0000",
				ValidThru:  "01/11",
				CVV:        "000",
				CardHolder: "name name",
			},
		},
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
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	cards, err := bankCardRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(cards), 2)
}

func TestSyncBankCardsServerError(t *testing.T) {
	client := FakeHTTPClient{
		Err: domain.ErrInvalidToken,
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
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestSyncBankCardsUnauthorized(t *testing.T) {
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
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
