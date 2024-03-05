// Package config содержит конфиг клиента
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const dbFileMode = 0600
const configName = "config.json"

type Config struct {
	DBFileMode    uint32
	APIVer        string
	DBFilePath    string        `json:"db_file_path"`
	ClientTimeout time.Duration `json:"db_client_timeout"`
	ServerURL     string        `json:"server_url"`
}

// New - Возвращает инстанс конфигурации сервера из файла
func New(root string) (*Config, error) {
	cfg := Config{
		DBFileMode:    dbFileMode,
		ClientTimeout: time.Duration(30) * time.Second, //nolint: gomnd
		APIVer:        "api/v1/",
		DBFilePath:    "user.db",
	}

	configPath := filepath.Join(root, configName)

	data, err := os.ReadFile(configPath) //nolint: gosec
	if err != nil {
		return &cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
