package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestUpdateTextSuccess(t *testing.T) {
	content := "my fancy content to update"
	textID := uuid.New()
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
		ID:      textID,
		Content: "old content",
	}
	err = textRepository.Create(userID, text)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"text",
		textID.String(),
		content,
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.Get(userID, textID)
	require.NoError(t, err)
	assert.Equal(t, txt.Content, content)
}

func TestUpdateTextBadRequest(t *testing.T) {
	textID := uuid.New()
	oldContent := "old content"
	client := FakeHTTPClient{
		Err: domain.ErrBadRequest,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	userID, err := createSession()
	require.NoError(t, err)

	text := domain.Text{
		ID:      textID,
		Content: oldContent,
	}
	err = textRepository.Create(userID, text)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"update",
		"text",
		textID.String(),
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	txt, err := textRepository.Get(userID, textID)
	require.NoError(t, err)
	assert.Equal(t, txt.Content, oldContent)
}

func TestUpdateTextUnauthorized(t *testing.T) {
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
		"text",
		"text_id",
		"content",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateTextNotFound(t *testing.T) {
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
		"text",
		uuid.NewString(),
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateTextInvalidValue(t *testing.T) {
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
		"text",
		"invalid value",
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
