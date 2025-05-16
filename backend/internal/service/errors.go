package service

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrSignToken     = errors.New("can't sign token")
	ErrParseToken    = errors.New("can't parse token")
	ErrParseExp      = errors.New("can't parse experience")

	ErrReservationNotFound = errors.New("reservation not found")
)
