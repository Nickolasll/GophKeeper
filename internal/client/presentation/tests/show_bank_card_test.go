package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestShowBankCardsSuccess(t *testing.T) {
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
		ID:         uuid.New(),
		Number:     "0000 0000 0000 0000",
		ValidThru:  "01/11",
		CVV:        "000",
		CardHolder: "",
	}
	err = bankCardRepository.Create(userID, &card)
	require.NoError(t, err)
	secondCard := domain.BankCard{
		ID:         uuid.New(),
		Number:     "1234 5678 1234 5678",
		ValidThru:  "11/11",
		CVV:        "222",
		CardHolder: "Name Name",
	}
	err = bankCardRepository.Create(userID, &secondCard)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"show",
		"bank-cards",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBankCardsNoContentSuccess(t *testing.T) {
	client := FakeHTTPClient{}

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
		"show",
		"bank-cards",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBankCardsNoToken(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"show",
		"bank-cards",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBankCardsExpiredToken(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	err = createExpiredSession()
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"show",
		"bank-cards",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
