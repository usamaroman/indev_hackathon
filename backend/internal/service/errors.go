package service

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrSignToken     = errors.New("can't sign token")
	ErrParseToken    = errors.New("can't parse token")
	ErrParseExp      = errors.New("can't parse experience")

	ErrBusinessExists    = errors.New("business already exists")
	ErrBusinessNotFound  = errors.New("business not found")
	ErrBusinessNotBelong = errors.New("business not belong to this user")

	ErrReservationNotFound = errors.New("reservation not found")

	ErrOverlap = errors.New("dates of reservation is overlap")
)
