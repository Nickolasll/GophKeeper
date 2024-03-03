// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
package domain

import "errors"

var ErrEntityNotFound = errors.New("entity not found")
var ErrUnauthorized = errors.New("invalid user credentials")
var ErrLoginConflict = errors.New("user with this login already exists")
var ErrBadRequest = errors.New("invalid input")
var ErrInvalidToken = errors.New("invalid token")
var ErrClientConnectionError = errors.New("http client connection error")
