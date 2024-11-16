package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/paudarco/referral-api/internal/errors"
	"github.com/paudarco/referral-api/internal/repository"
	"github.com/paudarco/referral-api/internal/storage"
	"github.com/paudarco/referral-api/models"
)

type ReferralService struct {
	db      repository.Referral
	storage *storage.Storage
}

func NewReferralService(db repository.Referral, s *storage.Storage) *ReferralService {
	return &ReferralService{db: db, storage: s}
}

func (s *ReferralService) CreateReferralCode(ctx context.Context, userID int, expires time.Time) (string, error) {
	if expires.Before(time.Now()) {
		return "", errors.ErrInvalidExpiration
	}

	// Generate unique code
	code, err := s.generateCode()
	if err != nil {
		return "", err
	}

	// Create new referral code
	referralCode := models.ReferralCode{
		ReferrerID: userID,
		Code:       code,
		ExpiresAt:  expires,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	// Save the referral code to the database
	err = s.db.SaveReferralCode(ctx, &referralCode)
	if err != nil {
		return "", err
	}

	// Storage the referral code
	s.storage.Set(referralCode.ID, referralCode, referralCode.ExpiresAt)

	return code, nil
}

func (s *ReferralService) GetReferrals(ctx context.Context, userId int) ([]models.User, error) {
	referrals, err := s.db.GetReferrals(ctx, userId)

	if err != nil {
		return nil, err
	}

	return referrals, nil
}

func (s *ReferralService) GetReferralCodeByEmail(ctx context.Context, email string) (string, error) {
	code, err := s.db.GetReferralById(ctx, email)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *ReferralService) DeleteReferralCode(ctx context.Context, userID int) error {
	id, err := s.db.DeactivateUserCodes(ctx, userID)
	if err != nil {
		return err
	}

	s.storage.DeleteUserCode(id)

	return nil
}

func (s *ReferralService) generateCode() (string, error) {
	bytes := make([]byte, 5)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
