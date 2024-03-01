package domain

import "errors"

var ErrLoginAlreadyInUse = errors.New("login already in use")
var ErrLoginOrPasswordIsInvalid = errors.New("login or password is invalid")
var ErrEntityNotFound = errors.New("entity not found")
