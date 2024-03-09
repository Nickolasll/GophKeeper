package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestSyncBinaryOverrideSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Binary{
			{
				ID:      uuid.New(),
				Content: []byte(uuid.NewString()),
			},
			{
				ID:      uuid.New(),
				Content: []byte(uuid.NewString()),
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

	binID := uuid.New()
	content := []byte("my fancy content")
	bin := domain.Binary{
		ID:      binID,
		Content: content,
	}
	err = binaryRepository.Create(userID, bin)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"sync",
		"binaries",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	bins, err := binaryRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(bins), 2)
	for _, v := range bins {
		assert.NotEqual(t, v.ID, bin)
		assert.NotEqual(t, v.Content, bin.Content)
	}
}

func TestSyncBinarySuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Binary{
			{
				ID:      uuid.New(),
				Content: []byte(uuid.NewString()),
			},
			{
				ID:      uuid.New(),
				Content: []byte(uuid.NewString()),
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
		"sync",
		"binaries",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	bins, err := binaryRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(bins), 2)
}

func TestSyncBinaryServerError(t *testing.T) {
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
		"binaries",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestSyncBinaryUnauthorized(t *testing.T) {
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
		"binaries",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
