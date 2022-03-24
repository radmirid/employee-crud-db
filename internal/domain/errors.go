package domain

import "errors"

var (
	ErrorEmployeeNotFound    = errors.New("employee not found")
	ErrorRefreshTokenExpired = errors.New("refresh token expired")
)
