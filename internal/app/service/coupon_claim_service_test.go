package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ijalalfrz/coupon-system/internal/app/dto"
	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCouponClaimService_ClaimCoupon(t *testing.T) {
	ctx := context.Background()

	// test helper
	claimCouponTest := func(
		name string,
		req dto.ClaimCouponRequest,
		setupMocks func(*MockCouponRepositoryClaimer, *MockCouponClaimRepositoryWriter),
		wantErr bool,
		expectedErr error,
	) func(t *testing.T) {
		return func(t *testing.T) {
			mockRedeemer := NewMockCouponRepositoryClaimer(t)
			mockWriter := NewMockCouponClaimRepositoryWriter(t)
			setupMocks(mockRedeemer, mockWriter)

			svc := NewCouponClaimService(mockRedeemer, mockWriter)
			err := svc.ClaimCoupon(ctx, req)

			if wantErr {
				assert.Error(t, err)
				if expectedErr != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		}
	}

	t.Run("success", claimCouponTest(
		"success",
		dto.ClaimCouponRequest{CouponName: "PROMO10", UserID: "user1"},
		func(mr *MockCouponRepositoryClaimer, mw *MockCouponClaimRepositoryWriter) {
			mr.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
					_ = txFunc(ctx, nil)
				}).
				Return(nil)

			mr.EXPECT().
				GetCouponForUpdateTx(mock.Anything, mock.Anything, "PROMO10").
				Return(model.Coupon{
					ID:              "1",
					Name:            "PROMO10",
					Amount:          10,
					RemainingAmount: 5,
				}, nil)

			mw.EXPECT().
				CreateCouponClaimTx(mock.Anything, mock.Anything, model.CouponClaim{
					CouponName: "PROMO10",
					UserID:     "user1",
				}).
				Return(nil)

			mr.EXPECT().
				SubstractRemainingAmountTx(mock.Anything, mock.Anything, "PROMO10", 1).
				Return(nil)
		},
		false,
		nil,
	))

	t.Run("no_stock", claimCouponTest(
		"no_stock",
		dto.ClaimCouponRequest{CouponName: "PROMO10", UserID: "user1"},
		func(mr *MockCouponRepositoryClaimer, mw *MockCouponClaimRepositoryWriter) {
			mr.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
					_ = txFunc(ctx, nil)
				}).
				Return(ErrCouponStockIsNotEnough) // WithTransaction returns the error from callback

			mr.EXPECT().
				GetCouponForUpdateTx(mock.Anything, mock.Anything, "PROMO10").
				Return(model.Coupon{
					ID:              "1",
					Name:            "PROMO10",
					Amount:          10,
					RemainingAmount: 0, // No stock
				}, nil)
		},
		true,
		ErrCouponStockIsNotEnough,
	))

	t.Run("already_claimed", claimCouponTest(
		"already_claimed",
		dto.ClaimCouponRequest{CouponName: "PROMO10", UserID: "user1"},
		func(mr *MockCouponRepositoryClaimer, mw *MockCouponClaimRepositoryWriter) {
			expectedErr := errors.New("coupon already claimed")

			mr.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
					_ = txFunc(ctx, nil)
				}).
				Return(expectedErr)

			mr.EXPECT().
				GetCouponForUpdateTx(mock.Anything, mock.Anything, "PROMO10").
				Return(model.Coupon{
					ID:              "1",
					Name:            "PROMO10",
					Amount:          10,
					RemainingAmount: 5,
				}, nil)

			mw.EXPECT().
				CreateCouponClaimTx(mock.Anything, mock.Anything, mock.Anything).
				Return(expectedErr)
		},
		true,
		errors.New("coupon already claimed"),
	))

	t.Run("db_error_get_coupon", claimCouponTest(
		"db_error_get_coupon",
		dto.ClaimCouponRequest{CouponName: "PROMO10", UserID: "user1"},
		func(mr *MockCouponRepositoryClaimer, mw *MockCouponClaimRepositoryWriter) {
			mr.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
					_ = txFunc(ctx, nil)
				}).
				Return(errors.New("db error"))

			mr.EXPECT().
				GetCouponForUpdateTx(mock.Anything, mock.Anything, "PROMO10").
				Return(model.Coupon{}, errors.New("db error"))
		},
		true,
		nil,
	))
}
