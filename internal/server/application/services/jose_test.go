package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJOSEIssue(t *testing.T) {
	tests := []struct {
		name string
		want uuid.UUID
	}{
		{
			name: "Issue and get claims",
			want: uuid.New(),
		},
	}
	for _, tt := range tests {
		raw := []byte("My secret keys")
		key, err := jwk.FromRaw(raw)
		require.NoError(t, err)
		joseService := JOSEService{
			TokenExp: time.Duration(60) * time.Second,
			JWKS:     key,
		}
		t.Run(tt.name, func(t *testing.T) {
			token, err := joseService.IssueToken(tt.want)
			require.NoError(t, err)
			userID, err := joseService.ParseUserID(token)
			require.NoError(t, err)
			assert.Equal(t, userID, tt.want)
		})
	}
}

func TestJOSEHash(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Issue and get claims random",
			want: uuid.NewString(),
		},
		{
			name: "Issue and get claims empty",
			want: "",
		},
	}
	for _, tt := range tests {
		raw := []byte("My secret keys")
		key, err := jwk.FromRaw(raw)
		require.NoError(t, err)
		joseService := JOSEService{
			TokenExp: time.Duration(60) * time.Second,
			JWKS:     key,
		}
		t.Run(tt.name, func(t *testing.T) {
			hash1 := joseService.Hash(tt.want)
			hash2 := joseService.Hash(tt.want)
			assert.NotEqual(t, hash1, hash2)
			ok := joseService.VerifyPassword(hash1, tt.want)
			assert.Equal(t, ok, true)
			ok = joseService.VerifyPassword(hash2, tt.want)
			assert.Equal(t, ok, true)
			ok = joseService.VerifyPassword(hash1, "wrong password")
			assert.Equal(t, ok, false)
			ok = joseService.VerifyPassword(hash2, "wrong password")
			assert.Equal(t, ok, false)
		})
	}
}
