package models

import "time"

// User представляет пользователя системы
type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// RegisterRequest представляет данные для регистрации
type RegisterRequest struct {
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	ReferralCode string `json:"referral_code,omitempty"`
}

// LoginRequest представляет данные для входа
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse представляет ответ при успешной аутентификации
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
