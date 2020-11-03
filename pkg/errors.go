package sso

import "errors"

var (
	UserOrPasswordIncorrect = errors.New("user or password incorrect")
	UserNotExists = errors.New("user not exists")
	TokenIncorrect = errors.New("token incorrect")
	RefreshTokenIncorrect = errors.New("refresh token incorrect")
)