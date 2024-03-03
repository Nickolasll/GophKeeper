// Package infrastructure содержит имплементацию репозиториев и клиентов
package infrastructure

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// HTTPClient - Имплементация клиента GophKeeper
type HTTPClient struct {
	client *resty.Client
}

// New - Конструктор нового инстанса клиента
func (HTTPClient) New(cert []byte, timeout time.Duration, baseURL string) HTTPClient {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
		RootCAs:       caCertPool,
		MinVersion:    tls.VersionTLS12,
	}

	client := resty.New().SetTLSClientConfig(tlsConfig).SetTimeout(timeout).SetBaseURL(baseURL)

	return HTTPClient{
		client: client,
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

// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
func (c HTTPClient) CreateText(session domain.Session, content string) (string, error) {
	resp, err := c.client.R().
		SetHeader("Content-Type", "plain/text").
		SetHeader("Authorization", session.Token).
		SetBody(content).
		Post("text/create")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == http.StatusCreated {
		return resp.Header().Get("Location"), nil
	}

	return "", domain.ErrClientConnectionError
}

// UpdateText - Обновляет существующий текст
func (c HTTPClient) UpdateText(session domain.Session, text domain.Text) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "plain/text").
		SetHeader("Authorization", session.Token).
		SetBody(text.Content).
		Post("text/" + text.ID)

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
