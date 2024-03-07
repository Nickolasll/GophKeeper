package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// HTTPClient - Имплементация клиента GophKeeper
type HTTPClient struct {
	client *resty.Client
	log    *logrus.Logger
}

// New - Конструктор нового инстанса клиента
func New(
	log *logrus.Logger,
	cert []byte,
	timeout time.Duration,
	baseURL string,
) *HTTPClient {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
		RootCAs:       caCertPool,
		MinVersion:    tls.VersionTLS13,
	}
	client := resty.New().
		SetTLSClientConfig(tlsConfig).
		SetTimeout(timeout).
		SetBaseURL(baseURL)

	return &HTTPClient{
		client: client,
		log:    log,
	}
}

// Login - Вход по логину и паролю, возвращает токен авторизации
func (c HTTPClient) Login(login, password string) (string, error) {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"login":    login,
			"password": password,
		}).Post("/auth/login")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return "", domain.ErrUnauthorized
	}

	if resp.StatusCode() == http.StatusOK {
		return resp.Header().Get("Authorization"), nil
	}

	return "", domain.ErrClientConnectionError
}

// Register - Регистрация по логину и паролю, возвращает токен авторизации
func (c HTTPClient) Register(login, password string) (string, error) {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"login":    login,
			"password": password,
		}).Post("/auth/register")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == http.StatusConflict {
		return "", domain.ErrLoginConflict
	}

	if resp.StatusCode() == http.StatusOK {
		return resp.Header().Get("Authorization"), nil
	}

	return "", domain.ErrClientConnectionError
}

func (c HTTPClient) create(authToken, uri, contentType string, body any) (string, error) {
	resp, err := c.client.R().
		SetHeader("Content-Type", contentType).
		SetHeader("Authorization", authToken).
		SetBody(body).
		Post(uri)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == http.StatusCreated {
		return resp.Header().Get("Location"), nil
	}

	return "", domain.ErrClientConnectionError
}

func (c HTTPClient) update(authToken, uri, contentType string, body any) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", contentType).
		SetHeader("Authorization", authToken).
		SetBody(body).
		Post(uri)

	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return domain.ErrEntityNotFound
	}
	if resp.StatusCode() == http.StatusBadRequest {
		return domain.ErrBadRequest
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}

	return domain.ErrClientConnectionError
}

// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
func (c HTTPClient) CreateText(session domain.Session, content string) (uuid.UUID, error) {
	var uid uuid.UUID
	id, err := c.create(session.Token, "text/create", "plain/text", content)

	if err != nil {
		return uid, err
	}

	uid, err = c.parseID(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// UpdateText - Обновляет существующий текст
func (c HTTPClient) UpdateText(session domain.Session, text domain.Text) error {
	err := c.update(session.Token, "text/"+text.ID.String(), "plain/text", text.Content)

	if err != nil {
		return err
	}

	return nil
}

// GetCerts - Возвращает публичный ключ для валидации и парсинга JWT
func (c HTTPClient) GetCerts() ([]byte, error) {
	resp, err := c.client.R().Get("/auth/certs")

	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode() == http.StatusOK {
		return resp.Body(), nil
	}

	return []byte{}, domain.ErrClientConnectionError
}

// CreateBinary - Создает бинарные данные, возвращает идентификатор ресурса от сервера
func (c HTTPClient) CreateBinary(session domain.Session, content []byte) (uuid.UUID, error) {
	var uid uuid.UUID
	id, err := c.create(session.Token, "binary/create", "multipart/form-data", content)

	if err != nil {
		return uid, err
	}

	uid, err = c.parseID(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// UpdateBinary - Обновляет существующие бинарные данные
func (c HTTPClient) UpdateBinary(session domain.Session, bin domain.Binary) error {
	err := c.update(session.Token, "binary/"+bin.ID.String(), "multipart/form-data", bin.Content)

	if err != nil {
		return err
	}

	return nil
}

func (c HTTPClient) parseID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// CreateCredentials - Создает пару логин и парль, возвращает идентификатор ресурса от сервера
func (c HTTPClient) CreateCredentials(session domain.Session, name, login, password string) (uuid.UUID, error) {
	var uid uuid.UUID
	payload, err := credentialsToJSON(name, login, password)
	if err != nil {
		return uid, err
	}
	id, err := c.create(session.Token, "credentials/create", "application/json", payload)

	if err != nil {
		return uid, err
	}

	uid, err = c.parseID(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// UpdateCredentials - Обновляет существующие логин и пароль
func (c HTTPClient) UpdateCredentials(session domain.Session, cred domain.Credentials) error {
	payload, err := credentialsToJSON(cred.Name, cred.Login, cred.Password)
	if err != nil {
		return err
	}
	err = c.update(session.Token, "credentials/"+cred.ID.String(), "application/json", payload)

	if err != nil {
		return err
	}

	return nil
}

// CreateBankCard - Создает банковскую карту, возвращает идентификатор ресурса от сервера
func (c HTTPClient) CreateBankCard(session domain.Session, number, validThru, cvv, cardHolder string) (uuid.UUID, error) {
	var uid uuid.UUID
	payload, err := bankCardToJSON(number, validThru, cvv, cardHolder)
	if err != nil {
		return uid, err
	}
	id, err := c.create(session.Token, "bank_card/create", "application/json", payload)

	if err != nil {
		return uid, err
	}

	uid, err = c.parseID(id)
	if err != nil {
		return uid, err
	}

	return uid, nil
}

// UpdateBankCard - Обновляет существующую банковскую карту
func (c HTTPClient) UpdateBankCard(session domain.Session, card *domain.BankCard) error {
	payload, err := bankCardToJSON(card.Number, card.ValidThru, card.CVV, card.CardHolder)
	if err != nil {
		return err
	}
	err = c.update(session.Token, "bank_card/"+card.ID.String(), "application/json", payload)

	if err != nil {
		return err
	}

	return nil
}
