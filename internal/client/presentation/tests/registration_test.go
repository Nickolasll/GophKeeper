package tests //nolint: dupl

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

func TestRegistrationSuccess(t *testing.T) {
	userID := uuid.NewString()
	token, err := issueToken(userID, time.Hour)
	require.NoError(t, err)
	client := FakeHTTPClient{
		Response: token,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"register",
		"test_login",
		"test_password",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)

	session, err := sessionRepository.Get()
	require.NoError(t, err)
	assert.Equal(t, session.Token, token)
	assert.Equal(t, session.UserID, userID)
}

func TestRegistrationConflict(t *testing.T) {
	client := FakeHTTPClient{
		Err: domain.ErrLoginConflict,
	}

	cmd, err := setup(client)
	require.NoError(t, err)
	defer func() {
		err = teardown()
		require.NoError(t, err)
	}()

	args := []string{
		"gophkeeper",
		"register",
		"test_login",
		"test_password",
	}

	err = cmd.Run(context.Background(), args)
	require.NoError(t, err)
}
