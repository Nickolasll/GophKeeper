package application

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCryptoInvalidSecretKey(t *testing.T) {
	cryptoService := CryptoService{
		SecretKey: []byte("w"),
	}
	_, err := cryptoService.Encrypt([]byte("test"))
	require.Error(t, err)
	_, err = cryptoService.Decrypt([]byte("test"))
	require.Error(t, err)
}
