package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func getClient() FakeHTTPClient {
	return FakeHTTPClient{
		SyncAllData: getAllResponse{
			Texts: []domain.Text{
				{
					ID:      uuid.New(),
					Content: uuid.NewString(),
				},
				{
					ID:      uuid.New(),
					Content: uuid.NewString(),
				},
			},
			Binaries: []domain.Binary{
				{
					ID:      uuid.New(),
					Content: []byte(uuid.NewString()),
				},
				{
					ID:      uuid.New(),
					Content: []byte(uuid.NewString()),
				},
			},
			Credentials: []domain.Credentials{
				{
					ID:       uuid.New(),
					Name:     "name",
					Login:    "login",
					Password: "password",
				},
				{
					ID:       uuid.New(),
					Name:     "name",
					Login:    "login",
					Password: "password",
				},
			},
			BankCards: []domain.BankCard{
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
		},
	}
}

func TestSyncAllOverrideSuccess(t *testing.T) {
	client := getClient()

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	text := domain.Text{
		ID:      uuid.New(),
		Content: "my fancy content",
	}
	err = textRepository.Create(userID, text)
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

	cred := domain.Credentials{
		ID:       uuid.New(),
		Name:     "my old name",
		Login:    "my old login",
		Password: "my old password",
	}
	err = credentialsRepository.Create(userID, &cred)
	require.NoError(t, err)

	bin := domain.Binary{
		ID:      uuid.New(),
		Content: []byte("my fancy content"),
	}
	err = binaryRepository.Create(userID, bin)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"sync",
		"all",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(txt), 2)

	bins, err := binaryRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(bins), 2)

	creds, err := credentialsRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(creds), 2)

	cards, err := bankCardRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(cards), 2)
}

func TestSyncAllSuccess(t *testing.T) {
	client := getClient()

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
		"sync",
		"all",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(txt), 2)

	bins, err := binaryRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(bins), 2)

	creds, err := credentialsRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(creds), 2)

	cards, err := bankCardRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(cards), 2)
}

func TestSyncAllServerError(t *testing.T) {
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
		"sync",
		"all",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestSyncAllUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"sync",
		"all",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
