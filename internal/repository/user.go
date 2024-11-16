package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paudarco/referral-api/internal/errors"
	"github.com/paudarco/referral-api/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.ErrDuplicateEmail
	}

	query := `
        INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `

	return r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
        SELECT id, name, email, password_hash, created_at
        FROM users
        WHERE email = $1
    `

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}
