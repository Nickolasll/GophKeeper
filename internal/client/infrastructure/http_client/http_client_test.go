package httpclient

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	log := logrus.New()
	cert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"SE"},
			Organization: []string{"Company Co."},
			CommonName:   "Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	client := New(log, cert.Raw, time.Minute, "http://test.url")
	require.NotNil(t, client)
}
