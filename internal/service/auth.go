package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/paudarco/referral-api/internal/config"
	"github.com/paudarco/referral-api/internal/errors"
	"github.com/paudarco/referral-api/internal/repository"
	"github.com/paudarco/referral-api/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repository.Authorization
	referralRepo repository.Referral
	jwtSecret    []byte
	tokenExpiry  time.Duration
}

func NewAuthService(repoAuth repository.Authorization, repoRef repository.Referral, cfg config.Config) *AuthService {
	return &AuthService{
		userRepo:     repoAuth,
		referralRepo: repoRef,
		jwtSecret:    []byte(cfg.JWT.Secret),
		tokenExpiry:  time.Duration(cfg.JWT.TTL) * time.Hour,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Check if referral code is valid and save it if present
	var codeId, referrerId int
	if req.ReferralCode != "" {
		codeId, referrerId, err = s.referralRepo.VerifyReferralCode(ctx, req.ReferralCode)
		if err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if req.ReferralCode != "" {
		err := s.referralRepo.SaveReferral(ctx, referrerId, user.ID, codeId)
		if err != nil {
			return nil, err
		}
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenExpiry).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Email,
	})

	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidSigningMethod
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return 0, errors.ErrInvalidClaims
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return 0, errors.ErrExpiredToken
	}

	return claims.UserID, nil

}
