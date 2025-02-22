package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestUpdateCredentialsSuccess(t *testing.T) {
	name := "new name"
	password := "new password"
	meta := "my meta"

	credID := uuid.New()
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
		ID:       credID,
		Name:     "old name",
		Login:    "old login",
		Password: "old password",
		Meta:     "",
	}
	err = credentialsRepository.Create(userID, &cred)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"credentials",
		"--meta",
		meta,
		"--name",
		name,
		"--password",
		password,
		credID.String(),
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	credObj, err := credentialsRepository.Get(userID, credID)
	require.NoError(t, err)
	assert.Equal(t, credObj.Name, name)
	assert.Equal(t, credObj.Login, "old login")
	assert.Equal(t, credObj.Password, password)
	assert.Equal(t, credObj.Meta, meta)
}

func TestUpdateCredentialsBadRequest(t *testing.T) {
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
		"credentials",
		uuid.NewString(),
	}
	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateCredentialsUnauthorized(t *testing.T) {
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
		"credentials",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateCredentialsNotFound(t *testing.T) {
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
		"credentials",
		"--name",
		"name",
		uuid.NewString(),
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateCredentialsInvalidValue(t *testing.T) {
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
		"credentials",
		"--name",
		"name",
		"invalid value",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
