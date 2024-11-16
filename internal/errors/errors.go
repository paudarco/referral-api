package errors

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user does not exist")
	ErrCodeNotFound         = errors.New("referral code not found")
	ErrCodeExpired          = errors.New("referral code expired")
	ErrCodeInactive         = errors.New("referral code inactive")
	ErrInvalidExpiration    = errors.New("invalid expiration")
	ErrDuplicateEmail       = errors.New("email already in use")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrExpiredToken         = errors.New("token expired")
	ErrUserIdNotFound       = errors.New("user id not found")
	ErrInvalidTypeId        = errors.New("invalid type user id")
)
