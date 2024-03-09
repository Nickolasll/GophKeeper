package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestSyncCredentialsOverrideSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Credentials{
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
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	cred := domain.Credentials{
		ID:       uuid.New(),
		Name:     "my old name",
		Login:    "my old login",
		Password: "my old password",
	}
	err = credentialsRepository.Create(userID, cred)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"credentials",
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	creds, err := credentialsRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(creds), 2)
	for _, v := range creds {
		assert.NotEqual(t, v.ID, cred.ID)
		assert.NotEqual(t, v.Name, cred.Name)
		assert.NotEqual(t, v.Login, cred.Login)
		assert.NotEqual(t, v.Password, cred.Password)
	}
}

func TestSyncCredentialsSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Credentials{
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
		"credentials",
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	creds, err := credentialsRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(creds), 2)
}

func TestSyncCredentialsServerError(t *testing.T) {
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
		"credentials",
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestSyncCredentialsUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"credentials",
		"sync",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
