package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestCreateCredentialsSuccess(t *testing.T) {
	name := "cred name"
	login := "login"
	password := "password"

	credID := uuid.New()
	client := FakeHTTPClient{
		Response: credID,
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
		"create",
		"credentials",
		name,
		login,
		password,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	cred, err := credentialsRepository.Get(userID, credID)
	require.NoError(t, err)
	assert.Equal(t, cred.Name, name)
	assert.Equal(t, cred.Login, login)
	assert.Equal(t, cred.Password, password)
}

func TestCreateCredentialsWithMetaSuccess(t *testing.T) {
	name := "cred name"
	login := "login"
	password := "password"
	meta := "meta"

	credID := uuid.New()
	client := FakeHTTPClient{
		Response: credID,
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
		"create",
		"credentials",
		"--meta",
		meta,
		name,
		login,
		password,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	cred, err := credentialsRepository.Get(userID, credID)
	require.NoError(t, err)
	assert.Equal(t, cred.Name, name)
	assert.Equal(t, cred.Login, login)
	assert.Equal(t, cred.Password, password)
}

func TestCreateCredentialsBadRequest(t *testing.T) {
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
		"create",
		"credentials",
		"name",
		"login",
		"password",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateCredentialsUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"create",
		"credentials",
		"name",
		"login",
		"password",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
