package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestShowCredentialsSuccess(t *testing.T) {
	client := FakeHTTPClient{}

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
		Name:     "www.example.com",
		Login:    "login",
		Password: "password",
	}
	err = credentialsRepository.Create(userID, &cred)
	require.NoError(t, err)
	secondCred := domain.Credentials{
		ID:       uuid.New(),
		Name:     "www.example2.com",
		Login:    "login2",
		Password: "password2",
	}
	err = credentialsRepository.Create(userID, &secondCred)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"show",
		"credentials",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowCredentialsNoContentSuccess(t *testing.T) {
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
		"credentials",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowCredentialNoToken(t *testing.T) {
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
		"credentials",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowCredentialsExpiredToken(t *testing.T) {
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
		"credentials",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
