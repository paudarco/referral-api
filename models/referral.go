package models

import "time"

// ReferralCode представляет реферальный код
type ReferralCode struct {
	ID         int       `json:"id" db:"id"`
	ReferrerID int       `json:"referrer_id" db:"referrer_id"`
	Code       string    `json:"code" db:"code"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CreateReferralCodeRequest представляет запрос на создание кода
type CreateReferralCodeRequest struct {
	ExpiresAt time.Time `json:"expires_at" validate:"required,gt=now"`
}

// ReferralInfo представляет информацию о реферале
type ReferralInfo struct {
	ReferralID int       `json:"referral_id"`
	Email      string    `json:"email"`
	JoinedAt   time.Time `json:"joined_at"`
}
