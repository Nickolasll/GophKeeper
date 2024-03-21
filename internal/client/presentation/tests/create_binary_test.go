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

func TestCreateBinarySuccess(t *testing.T) {
	binID := uuid.New()
	client := FakeHTTPClient{
		Response: binID,
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
		"binary",
		"./binary_file_for_test",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	bin, err := binaryRepository.Get(userID, binID)
	require.NoError(t, err)
	content, err := os.ReadFile("./binary_file_for_test")
	require.NoError(t, err)
	assert.Equal(t, bin.Content, content)
}

func TestCreateBinaryBadRequest(t *testing.T) {
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
		"binary",
		"./binary_file_for_test",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateBinaryUnauthorized(t *testing.T) {
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
		"binary",
		"./binary_file_for_test",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}

func TestCreateBinaryNoFile(t *testing.T) {
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
		"create",
		"binary",
		"invalid",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
