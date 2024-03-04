package infrastructure

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
		MinVersion:    tls.VersionTLS13,
	}

	client := HTTPClient{}.New(tlsConfig, time.Minute, "http://test.url")
	require.NotNil(t, client)
}
