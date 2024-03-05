package tests

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestUpdateBinarySuccess(t *testing.T) {
	binID := uuid.New()
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
		ID:      binID,
		Content: []byte("old content"),
	}
	err = binaryRepository.Create(userID, bin)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"binary",
		"update",
		binID.String(),
		"./binary_file_for_test",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	binObj, err := binaryRepository.Get(userID, binID)
	require.NoError(t, err)
	content, err := os.ReadFile("./binary_file_for_test")
	require.NoError(t, err)
	assert.Equal(t, binObj.Content, content)
}

func TestUpdateBinaryBadRequest(t *testing.T) {
	binID := uuid.New()
	oldContent := []byte("old content")
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

	bin := domain.Binary{
		ID:      binID,
		Content: oldContent,
	}
	err = binaryRepository.Create(userID, bin)
	require.NoError(t, err)

	args := []string{
		"gophkeeper",
		"binary",
		"update",
		binID.String(),
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	binObj, err := binaryRepository.Get(userID, binID)
	require.NoError(t, err)
	assert.Equal(t, binObj.Content, oldContent)
}

func TestUpdateBinaryUnauthorized(t *testing.T) {
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
		"update",
		"text_id",
		"content",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBinaryNotFound(t *testing.T) {
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
		"binary",
		"update",
		uuid.NewString(),
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestUpdateBinaryInvalidUUID(t *testing.T) {
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
		"binary",
		"update",
		"invalid value",
		"",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
