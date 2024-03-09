package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestSyncTextOverrideSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Text{
			{
				ID:      uuid.New(),
				Content: uuid.NewString(),
			},
			{
				ID:      uuid.New(),
				Content: uuid.NewString(),
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

	textID := uuid.New()
	content := "my fancy content"
	text := domain.Text{
		ID:      textID,
		Content: content,
	}
	err = textRepository.Create(userID, text)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"sync",
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(txt), 2)
	for _, v := range txt {
		assert.NotEqual(t, v.ID, textID)
		assert.NotEqual(t, v.Content, text.Content)
	}
}

func TestSyncTextSuccess(t *testing.T) {
	client := FakeHTTPClient{
		Response: []domain.Text{
			{
				ID:      uuid.New(),
				Content: uuid.NewString(),
			},
			{
				ID:      uuid.New(),
				Content: uuid.NewString(),
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.GetAll(userID)
	require.NoError(t, err)
	assert.Equal(t, len(txt), 2)
}

func TestSyncTextServerError(t *testing.T) {
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestSyncTextUnauthorized(t *testing.T) {
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
		"texts",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
