package httpclient

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

const registerPath = "/auth/register"
const loginPath = "/auth/login"
const textPath = "/text/"
const textCreatePath = textPath + "create"
const textAllPath = textPath + "all"

func newClient(url string) *HTTPClient {
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

	return New(log, cert.Raw, time.Second, url)
}

func newSession() domain.Session {
	return domain.Session{
		UserID: uuid.New(),
		Token:  "tokenValue",
	}
}

func TestNewClient(t *testing.T) {
	client := newClient("http://test.url")
	require.NotNil(t, client)
}

func TestGetCertsInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/certs" {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"message": "Something went wrong!"}`)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	_, err := client.GetCerts()
	require.Error(t, err)
}

func TestGetCertsWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")

	_, err := client.GetCerts()
	require.Error(t, err)
}

func TestGetCertsSuccess(t *testing.T) {
	publicKey := []byte(`{"alg": "RS256", "kty": "RSA"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/certs" {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(publicKey); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	data, err := client.GetCerts()
	require.NoError(t, err)
	assert.Equal(t, publicKey, data)
}

func TestRegisterSuccess(t *testing.T) {
	tokenValue := "tokenValue"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == registerPath {
			w.Header().Set("Authorization", tokenValue)
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	token, err := client.Register("login", "password")
	require.NoError(t, err)
	assert.Equal(t, token, tokenValue)
}

func TestRegisterInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == registerPath {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"message": "Something went wrong!"}`)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	_, err := client.Register("login", "password")
	require.Error(t, err)
}

func TestRegisterLoginIsTaken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == registerPath {
			w.WriteHeader(http.StatusConflict)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	_, err := client.Register("login", "password")
	require.Error(t, err)
}

func TestRegisterWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")

	_, err := client.Register("login", "password")
	require.Error(t, err)
}

func TestLoginSuccess(t *testing.T) {
	tokenValue := "tokenValue"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == loginPath {
			w.Header().Set("Authorization", tokenValue)
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	token, err := client.Login("login", "password")
	require.NoError(t, err)
	assert.Equal(t, token, tokenValue)
}

func TestLoginWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")

	_, err := client.Login("login", "password")
	require.Error(t, err)
}

func TestLoginInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == loginPath {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	_, err := client.Login("login", "password")
	require.Error(t, err)
}

func TestLoginInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == loginPath {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"message": "Something went wrong!"}`)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)

	_, err := client.Login("login", "password")
	require.Error(t, err)
}

func TestCreateTextWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()

	_, err := client.CreateText(session, "content")
	require.Error(t, err)
}

func TestCreateInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textCreatePath {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"message": "Something went wrong!"}`)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.CreateText(session, "content")
	require.Error(t, err)
}

func TestCreateTextSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textCreatePath {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	uid, err := client.CreateText(session, "content")
	require.NoError(t, err)
	assert.Equal(t, uid, id)
}

func TestCreateTextInvalidLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textCreatePath {
			w.Header().Set("Location", "invalid")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.CreateText(session, "content")
	require.Error(t, err)
}

func TestCreateBinarySuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/binary/create" {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	uid, err := client.CreateBinary(session, []byte("content"))
	require.NoError(t, err)
	assert.Equal(t, uid, id)
}

func TestCreateBinaryInvalidLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/binary/create" {
			w.Header().Set("Location", "invalid")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.CreateBinary(session, []byte("content"))
	require.Error(t, err)
}

func TestCreateBinaryWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()

	_, err := client.CreateBinary(session, []byte("content"))
	require.Error(t, err)
}

func TestCreateCredentialsSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/credentials/create" {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	uid, err := client.CreateCredentials(session, "name", "login", "password")
	require.NoError(t, err)
	assert.Equal(t, uid, id)
}

func TestCreateCredentialsInvalidLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/credentials/create" {
			w.Header().Set("Location", "invalid")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.CreateCredentials(session, "name", "login", "password")
	require.Error(t, err)
}

func TestCreateCredentialsWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()

	_, err := client.CreateCredentials(session, "name", "login", "password")
	require.Error(t, err)
}

func TestCreateBankCardSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bank_card/create" {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	uid, err := client.CreateBankCard(session, "number", "valid_thru", "cvv", "card_holder")
	require.NoError(t, err)
	assert.Equal(t, uid, id)
}

func TestCreateBankCardInvalidLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bank_card/create" {
			w.Header().Set("Location", "invalid")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.CreateBankCard(session, "number", "valid_thru", "cvv", "card_holder")
	require.Error(t, err)
}

func TestCreateBankCardWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()

	_, err := client.CreateBankCard(session, "number", "valid_thru", "cvv", "card_holder")
	require.Error(t, err)
}

func TestUpdateNotFound(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textPath+id.String() {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	text := domain.Text{
		ID:      id,
		Content: "content",
	}

	err := client.UpdateText(session, text)
	require.Error(t, err)
}

func TestUpdateBadRequest(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textPath+id.String() {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	text := domain.Text{
		ID:      id,
		Content: "content",
	}

	err := client.UpdateText(session, text)
	require.Error(t, err)
}

func TestUpdateInvalidLocation(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textPath+id.String() {
			w.Header().Set("Location", "invalid")
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()
	text := domain.Text{
		ID:      id,
		Content: "content",
	}

	client := newClient(server.URL)
	session := newSession()

	err := client.UpdateText(session, text)
	require.Error(t, err)
}

func TestUpdateTextSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textPath+id.String() {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	text := domain.Text{
		ID:      id,
		Content: "content",
	}

	err := client.UpdateText(session, text)
	require.NoError(t, err)
}

func TestUpdateTextWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()
	text := domain.Text{
		ID:      uuid.New(),
		Content: "content",
	}

	err := client.UpdateText(session, text)
	require.Error(t, err)
}

func TestUpdateBinarySuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/binary/"+id.String() {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	bin := domain.Binary{
		ID:      id,
		Content: []byte("content"),
	}

	err := client.UpdateBinary(session, bin)
	require.NoError(t, err)
}

func TestUpdateBinaryWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()
	bin := domain.Binary{
		ID:      uuid.New(),
		Content: []byte("content"),
	}

	err := client.UpdateBinary(session, bin)
	require.Error(t, err)
}

func TestUpdateCredentialsSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/credentials/"+id.String() {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	cred := domain.Credentials{
		ID:       id,
		Name:     "name",
		Login:    "login",
		Password: "password",
	}

	err := client.UpdateCredentials(session, cred)
	require.NoError(t, err)
}

func TestUpdateCredentialsWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()
	cred := domain.Credentials{
		ID:       uuid.New(),
		Name:     "name",
		Login:    "login",
		Password: "password",
	}

	err := client.UpdateCredentials(session, cred)
	require.Error(t, err)
}

func TestUpdateBankCardSuccess(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bank_card/"+id.String() {
			w.Header().Set("Location", id.String())
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()
	card := domain.BankCard{
		ID:         id,
		Number:     "number",
		ValidThru:  "valid_thru",
		CVV:        "cvv",
		CardHolder: "card_holder",
	}

	err := client.UpdateBankCard(session, &card)
	require.NoError(t, err)
}

func TestUpdateBankCardWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()
	card := domain.BankCard{
		ID:         uuid.New(),
		Number:     "number",
		ValidThru:  "valid_thru",
		CVV:        "cvv",
		CardHolder: "card_holder",
	}

	err := client.UpdateBankCard(session, &card)
	require.Error(t, err)
}

func TestGetAllTextsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textAllPath {
			response := getAllTextsResponse{}
			response.Data.Texts = []domain.Text{
				{
					ID:      uuid.New(),
					Content: "content",
				},
				{
					ID:      uuid.New(),
					Content: "content",
				},
			}
			respData, err := json.Marshal(response)
			if err != nil {
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err = w.Write(respData); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	data, err := client.GetAllTexts(session)
	require.NoError(t, err)
	assert.Equal(t, len(data), 2)
}

func TestGetAllTextsWrongURL(t *testing.T) {
	client := newClient("wrongurl.com")
	session := newSession()

	_, err := client.GetAllTexts(session)
	require.Error(t, err)
}

func TestGetAllTextsInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textAllPath {
			response := errorResponse{Message: "error :("}
			respData, err := json.Marshal(response)
			if err != nil {
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			if _, err = w.Write(respData); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.GetAllTexts(session)
	require.Error(t, err)
}

func TestGetAllTextsBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == textAllPath {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := newClient(server.URL)
	session := newSession()

	_, err := client.GetAllTexts(session)
	require.Error(t, err)
}
