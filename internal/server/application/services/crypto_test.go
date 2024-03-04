package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Проверяем, что шифрование одного и того же сообщения
// каждый раз выдает разный результат
// При этом каждый раз результат корректно расшифровывается
func TestCrypto(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Encrypt Decrypt",
			want: "My test message",
		},
		{
			name: "Encrypt Decrypt Empty",
			want: "",
		},
		{
			name: "Encrypt Decrypt UUID",
			want: uuid.NewString(),
		},
		{
			name: "Encrypt Decrypt Short",
			want: "1",
		},
		{
			name: "Encrypt Decrypt Long",
			want: `Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
				  Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when 
				  an unknown printer took a galley of type and scrambled it to make a type specimen book.
				  It has survived not only five centuries, but also the leap into electronic typesetting, 
			 	  remaining essentially unchanged. It was popularized in the 1960s with the release of Letraset sheets 
			 	  containing Lorem Ipsum passages, and more recently with desktop publishing software like 
			 	  Aldus PageMaker including versions of Lorem Ipsum.`,
		},
	}
	for _, tt := range tests {
		cryptoService := CryptoService{
			SecretKey: []byte("1234567812345678"),
		}
		t.Run(tt.name, func(t *testing.T) {
			encr1, err := cryptoService.Encrypt([]byte(tt.want))
			require.NoError(t, err)
			encr2, err := cryptoService.Encrypt([]byte(tt.want))
			require.NoError(t, err)
			assert.NotEqual(t, encr1, encr2)
			decr1, err := cryptoService.Decrypt(encr1)
			require.NoError(t, err)
			assert.Equal(t, string(decr1), tt.want)
			decr2, err := cryptoService.Decrypt(encr2)
			require.NoError(t, err)
			assert.Equal(t, string(decr2), tt.want)
		})
	}
}

func TestCryptoInvalidSecretKey(t *testing.T) {
	cryptoService := CryptoService{
		SecretKey: []byte("w"),
	}
	data := []byte("test")
	_, err := cryptoService.Encrypt(data)
	require.Error(t, err)
	_, err = cryptoService.Decrypt(data)
	require.Error(t, err)
}
