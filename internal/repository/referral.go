package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paudarco/referral-api/models"
)

type ReferralRepository struct {
	db *pgxpool.Pool
}

func NewReferralRepository(db *pgxpool.Pool) *ReferralRepository {
	return &ReferralRepository{db: db}
}

// Insert new referral code into the database.
func (r *ReferralRepository) SaveReferralCode(ctx context.Context, code *models.ReferralCode) error {
	if _, err := r.DeactivateUserCodes(ctx, code.ReferrerID); err != nil {
		return err
	}

	query := `
        INSERT INTO referral_codes (referrer_id, code, expires_at, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `

	return r.db.QueryRow(
		ctx,
		query,
		code.ReferrerID,
		code.Code,
		code.ExpiresAt,
		code.IsActive,
	).Scan(&code.ID, &code.CreatedAt)
}

// Saving referral to the database
func (r *ReferralRepository) SaveReferral(ctx context.Context, referrerId, referralId, codeId int) error {

	query := "INSERT INTO referrals (referrer_id, referral_id, referral_code_id) VALUES ($1, $2, $3)"

	_, err := r.db.Exec(ctx, query, referrerId, referralId, codeId)
	if err != nil {
		return err
	}

	return nil
}

// Returns a list of all user's referrals
func (r *ReferralRepository) GetReferrals(ctx context.Context, userID int) ([]models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	var usersId []int
	query := "SELECT referral_id FROM referrals WHERE referrer_id = $1"

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		usersId = append(usersId, id)
	}

	query = "SELECT name, email, created_at FROM users WHERE id = $1"

	for id := range usersId {
		var user models.User
		err = tx.QueryRow(ctx, query, id).Scan(&user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		user.ID = id
		users = append(users, user)
	}

	return users, tx.Commit(ctx)
}

// Returns the active referral code associated with the given user's email
func (r *ReferralRepository) GetReferralById(ctx context.Context, email string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}

	// Get the user ID from the database.
	query := "SELECT id FROM users WHERE email = $1"

	var userId int
	err = tx.QueryRow(ctx, query, email).Scan(&userId)
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}

	query = "SELECT code FROM referral_codes WHERE referrer_id = $1 AND is_active = true"

	var code string
	err = tx.QueryRow(ctx, query, userId).Scan(&code)
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}

	return code, tx.Commit(ctx)
}

// Deactivates all active referral codes associated with the given userID
func (r *ReferralRepository) DeactivateUserCodes(ctx context.Context, userID int) (int, error) {
	query := `
        UPDATE referral_codes
        SET is_active = false
        WHERE referrer_id = $1
		RETURNING id
    `

	row := r.db.QueryRow(ctx, query, userID)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReferralRepository) VerifyReferralCode(ctx context.Context, code string) (int, int, error) {
	query := "SELECT id, referrer_id FROM referral_codes WHERE code = $1 AND is_active = true"

	var codeId, referrerId int
	err := r.db.QueryRow(ctx, query, code).Scan(&codeId, &referrerId)
	if err != nil {
		return 0, 0, err
	}

	return codeId, referrerId, nil
}
