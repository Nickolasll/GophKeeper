package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestCreateTextSuccess(t *testing.T) {
	content := "my fancy content for test"
	textID := uuid.New()
	client := FakeHTTPClient{
		Response: textID,
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
		"text",
		"create",
		content,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	text, err := textRepository.Get(userID, textID)
	require.NoError(t, err)
	assert.Equal(t, text.Content, content)
}

func TestCreateTextBadRequest(t *testing.T) {
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
		"text",
		"create",
		"content",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateTextUnauthorized(t *testing.T) {
	client := FakeHTTPClient{}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"text",
		"create",
		"content",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
