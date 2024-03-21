package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestShowTextSuccess(t *testing.T) {
	client := FakeHTTPClient{}

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
		Content: "old content",
	}
	err = textRepository.Create(userID, text)
	require.NoError(t, err)
	secondText := domain.Text{
		ID:      uuid.New(),
		Content: "second text",
	}
	err = textRepository.Create(userID, secondText)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"show",
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowTextNoContentSuccess(t *testing.T) {
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowTextNoToken(t *testing.T) {
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestShowTextExpiredToken(t *testing.T) {
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
