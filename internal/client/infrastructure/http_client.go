// Package infrastructure содержит имплементацию репозиториев и клиентов
package infrastructure

import (
	"crypto/tls"
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
func (HTTPClient) New(tlsConfig *tls.Config, timeout time.Duration, baseURL string) HTTPClient {

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
func (c HTTPClient) CreateText(session domain.Session, content string) (string, error) {
	id, err := c.create(session.Token, "text/create", "plain/text", content)

	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateText - Обновляет существующий текст
func (c HTTPClient) UpdateText(session domain.Session, text domain.Text) error {
	err := c.update(session.Token, "text/"+text.ID, "plain/text", text.Content)

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
func (c HTTPClient) CreateBinary(session domain.Session, content []byte) (string, error) {
	id, err := c.create(session.Token, "binary/create", "multipart/form-data", content)

	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateBinary - Обновляет существующие бинарные данные
func (c HTTPClient) UpdateBinary(session domain.Session, bin domain.Binary) error {
	err := c.update(session.Token, "binary/"+bin.ID, "multipart/form-data", bin.Content)

	if err != nil {
		return err
	}

	return nil
}
