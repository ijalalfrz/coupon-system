package service

import (
	"context"
	"fmt"

	"github.com/ijalalfrz/coupon-system/internal/app/dto"
	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/jmoiron/sqlx"
)

type CouponRepositoryClaimer interface {
	WithTransaction(ctx context.Context,
		txFunc func(context.Context, *sqlx.Tx) error,
	) error
	SubstractRemainingAmountTx(ctx context.Context, tx *sqlx.Tx, couponName string, amount int) error
	GetCouponForUpdateTx(ctx context.Context, tx *sqlx.Tx, couponName string) (model.Coupon, error)
}

type CouponClaimRepositoryWriter interface {
	CreateCouponClaimTx(ctx context.Context, tx *sqlx.Tx, couponClaim model.CouponClaim) error
}

type CouponClaimService struct {
	couponRepositoryClaimer     CouponRepositoryClaimer
	couponClaimRepositoryWriter CouponClaimRepositoryWriter
}

func NewCouponClaimService(couponRepositoryClaimer CouponRepositoryClaimer,
	couponClaimRepositoryWriter CouponClaimRepositoryWriter) *CouponClaimService {
	return &CouponClaimService{
		couponRepositoryClaimer:     couponRepositoryClaimer,
		couponClaimRepositoryWriter: couponClaimRepositoryWriter,
	}
}

func (s *CouponClaimService) ClaimCoupon(ctx context.Context, req dto.ClaimCouponRequest) error {
	couponClaim := model.CouponClaim{
		CouponName: req.CouponName,
		UserID:     req.UserID,
	}

	err := s.couponRepositoryClaimer.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		coupon, err := s.couponRepositoryClaimer.GetCouponForUpdateTx(ctx, tx, couponClaim.CouponName)
		if err != nil {
			return fmt.Errorf("failed to get coupon for update: %w", err)
		}

		if coupon.RemainingAmount <= 0 {
			return ErrCouponStockIsNotEnough
		}

		err = s.couponClaimRepositoryWriter.CreateCouponClaimTx(ctx, tx, couponClaim)
		if err != nil {
			return fmt.Errorf("failed to create coupon claim: %w", err)
		}

		err = s.couponRepositoryClaimer.SubstractRemainingAmountTx(ctx, tx, couponClaim.CouponName, 1)
		if err != nil {
			return fmt.Errorf("failed to substract remaining amount: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to claim coupon: %w", err)
	}

	return nil
}
