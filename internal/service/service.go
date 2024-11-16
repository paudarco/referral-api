package service

import (
	"context"
	"time"

	"github.com/paudarco/referral-api/internal/config"
	"github.com/paudarco/referral-api/internal/repository"
	"github.com/paudarco/referral-api/internal/storage"
	"github.com/paudarco/referral-api/models"
)

type Authorization interface {
	Register(ctx context.Context, user *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, user *models.LoginRequest) (*models.AuthResponse, error)
	ValidateToken(token string) (int, error)
}

type Referral interface {
	CreateReferralCode(c context.Context, userID int, expiresAt time.Time) (string, error)
	GetReferrals(ctx context.Context, userId int) ([]models.User, error)
	GetReferralCodeByEmail(ctx context.Context, email string) (string, error)
	DeleteReferralCode(ctx context.Context, userID int) error
}

type Service struct {
	Authorization
	Referral
}

func NewService(repos *repository.Repository, cfg config.Config, storage *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(repos, repos, cfg),
		Referral:      NewReferralService(repos, storage),
	}
}
