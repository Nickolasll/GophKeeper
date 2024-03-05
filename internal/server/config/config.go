// Package config содержит конфиг сервера
package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Config - конфигурация сервера
type Config struct {
	// Addr - Адрес сервера
	Addr string `env:"ADDR, default=localhost:8080"`
	// DBTimeOut - Таймаут операций бд
	DBTimeOut time.Duration `env:"DB_TIMEOUT, default=15s"`
	// JWTExpiration - Время жизни JWT
	JWTExpiration time.Duration `env:"JWT_EXPIRATION, default=600s"`
	// RawJWK - Json web keys в сериализованном формате
	RawJWK []byte `env:"RAW_JWK, default=My secret keys"`
	// PostgresURL - Путь до базы postgres
	PostgresURL string `env:"POSTGRES_URL, default=postgresql://admin:admin@localhost:5432/gophkeeper"`
	// CryptoSecret - Приватный ключ для шифрования данных
	CryptoSecret []byte `env:"CRYPTO_SECRET, default=1234567812345678"`
	// ReadHeaderTimeout - Таймаут чтения заголовков
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT, default=2s"`
	// X509CertPath - Путь до сертификата x509
	X509CertPath string `env:"X509_CERT_PATH, default=server.crt"`
	// X509KeyPath - Путь до ключа tls
	TLSKeyPath string `env:"TLS_KEY_PATH, default=server.key"`
}

// New - Возвращает инстанс конфигурации сервера из переменных окружения
func New() (*Config, error) {
	var cfg Config
	ctx := context.Background()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
