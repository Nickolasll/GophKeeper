package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestShowBinarySuccess(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	bin := domain.Binary{
		ID:      uuid.New(),
		Content: []byte("old content"),
	}
	err = binaryRepository.Create(userID, bin)
	require.NoError(t, err)
	secondBin := domain.Binary{
		ID:      uuid.New(),
		Content: []byte("second text"),
	}
	err = binaryRepository.Create(userID, secondBin)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"binary",
		"show",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBinaryNoContentSuccess(t *testing.T) {
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
		"binary",
		"show",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBinaryNoToken(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"binary",
		"show",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowBinaryExpiredToken(t *testing.T) {
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
		"binary",
		"show",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
