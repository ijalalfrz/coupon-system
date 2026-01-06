package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/ijalalfrz/coupon-system/internal/pkg/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type CouponRepository struct {
	clientDB *db.DB
	transactable
}

func NewCouponRepository(clientDB *db.DB) *CouponRepository {
	return &CouponRepository{
		clientDB: clientDB,
		transactable: transactable{
			db: clientDB.DB,
		},
	}
}

func (r *CouponRepository) CreateTx(ctx context.Context,
	tx *sqlx.Tx,
	coupon model.Coupon,
) error {
	query := "INSERT INTO coupons(name,amount,remaining_amount) VALUES($1,$2,$3)"

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, coupon.Name, coupon.Amount, coupon.RemainingAmount)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrRecordAlreadyExists
		}

		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil

}

func (r *CouponRepository) SubstractRemainingAmountTx(ctx context.Context,
	tx *sqlx.Tx,
	couponName string,
	amount int,
) error {
	query := "UPDATE coupons SET remaining_amount = remaining_amount - $1 WHERE name = $2 AND remaining_amount > 0"

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, amount, couponName)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("coupon not found or remaining amount is not enough")
	}

	return nil
}

func (r *CouponRepository) GetCouponForUpdateTx(ctx context.Context,
	tx *sqlx.Tx,
	couponName string,
) (model.Coupon, error) {
	query := "SELECT id,name,amount,remaining_amount FROM coupons WHERE name = $1 FOR UPDATE"

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return model.Coupon{}, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var coupon model.Coupon
	err = stmt.QueryRowContext(ctx, couponName).
		Scan(&coupon.ID, &coupon.Name, &coupon.Amount, &coupon.RemainingAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Coupon{}, ErrRecordNotFound
		}

		return model.Coupon{}, fmt.Errorf("failed to query row: %w", err)
	}

	return coupon, nil
}

func (r *CouponRepository) GetCouponByName(ctx context.Context,
	couponName string,
) (model.Coupon, error) {
	query := "SELECT id,name,amount,remaining_amount FROM coupons WHERE name = $1"

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return model.Coupon{}, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var coupon model.Coupon
	err = stmt.QueryRowContext(ctx, couponName).
		Scan(&coupon.ID, &coupon.Name, &coupon.Amount, &coupon.RemainingAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Coupon{}, ErrRecordNotFound
		}

		return model.Coupon{}, fmt.Errorf("failed to query row: %w", err)
	}

	return coupon, nil
}
