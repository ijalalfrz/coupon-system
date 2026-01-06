package service

import (
	"context"
	"fmt"

	"github.com/ijalalfrz/coupon-system/internal/app/dto"
	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/jmoiron/sqlx"
)

// CouponRepositoryWriter interface for coupon repository writer
type CouponRepositoryWriter interface {
	WithTransaction(ctx context.Context,
		txFunc func(context.Context, *sqlx.Tx) error,
	) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, coupon model.Coupon) error
}

// CouponRepositoryReader interface for coupon repository reader
type CouponRepositoryReader interface {
	GetCouponByName(ctx context.Context, couponName string) (model.Coupon, error)
}

// CouponClaimRepositoryReader interface for coupon claim repository reader
type CouponClaimRepositoryReader interface {
	GetAllCouponClaimByCouponName(ctx context.Context, couponName string) ([]string, error)
}

// CouponService struct for coupon service
type CouponService struct {
	couponRepositoryWriter      CouponRepositoryWriter
	couponRepositoryReader      CouponRepositoryReader
	couponClaimRepositoryReader CouponClaimRepositoryReader
}

func NewCouponService(couponRepositoryWriter CouponRepositoryWriter,
	couponRepositoryReader CouponRepositoryReader,
	couponClaimRepositoryReader CouponClaimRepositoryReader) *CouponService {
	return &CouponService{
		couponRepositoryWriter:      couponRepositoryWriter,
		couponRepositoryReader:      couponRepositoryReader,
		couponClaimRepositoryReader: couponClaimRepositoryReader,
	}
}

func (s *CouponService) CreateCoupon(ctx context.Context, req dto.CreateCouponRequest) error {
	coupon := model.Coupon{
		Name:            req.Name,
		Amount:          req.Amount,
		RemainingAmount: req.Amount,
	}

	err := s.couponRepositoryWriter.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.couponRepositoryWriter.CreateTx(ctx, tx, coupon)
	})

	if err != nil {
		return fmt.Errorf("failed to create coupon: %w", err)
	}

	return nil
}

func (s *CouponService) GetCouponByCouponName(ctx context.Context, req dto.GetCouponRequest) (dto.GetCouponResponse, error) {
	coupon, err := s.couponRepositoryReader.GetCouponByName(ctx, req.Name)
	if err != nil {
		return dto.GetCouponResponse{}, fmt.Errorf("failed to get coupon by coupon name: %w", err)
	}

	couponClaimUserIDs, err := s.couponClaimRepositoryReader.GetAllCouponClaimByCouponName(ctx, req.Name)
	if err != nil {
		return dto.GetCouponResponse{}, fmt.Errorf("failed to get coupon claim user ids: %w", err)
	}

	return dto.GetCouponResponse{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: coupon.RemainingAmount,
		ClaimedBy:       couponClaimUserIDs,
	}, nil
}
