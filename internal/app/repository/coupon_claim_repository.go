package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/ijalalfrz/coupon-system/internal/pkg/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type CouponClaimRepository struct {
	clientDB *db.DB
	transactable
}

func NewCouponClaimRepository(clientDB *db.DB) *CouponClaimRepository {
	return &CouponClaimRepository{
		clientDB: clientDB,
		transactable: transactable{
			db: clientDB.DB,
		},
	}
}

func (r *CouponClaimRepository) CreateCouponClaimTx(ctx context.Context,
	tx *sqlx.Tx,
	couponClaim model.CouponClaim,
) error {
	query := "INSERT INTO coupon_claims (coupon_name, user_id) VALUES ($1, $2)"

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, couponClaim.CouponName, couponClaim.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrCouponAlreadyClaimed
		}

		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (r *CouponClaimRepository) GetAllCouponClaimByCouponName(ctx context.Context,
	couponName string,
) ([]string, error) {
	query := "SELECT user_id FROM coupon_claims WHERE coupon_name = $1"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var userIDs []string

	rows, err := stmt.QueryContext(ctx, couponName)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}
