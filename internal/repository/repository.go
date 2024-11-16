package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paudarco/referral-api/models"
)

type Authorization interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type Referral interface {
	SaveReferralCode(ctx context.Context, code *models.ReferralCode) error
	GetReferrals(ctx context.Context, userID int) ([]models.User, error)
	GetReferralById(ctx context.Context, email string) (string, error)
	DeactivateUserCodes(ctx context.Context, userID int) (int, error)
	VerifyReferralCode(ctx context.Context, code string) (int, int, error)
	SaveReferral(ctx context.Context, referrerId, referralId, codeId int) error
}

type Repository struct {
	Authorization
	Referral
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewUserRepository(db),
		Referral:      NewReferralRepository(db),
	}
}
