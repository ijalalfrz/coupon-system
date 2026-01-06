package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ijalalfrz/coupon-system/internal/app/dto"
	"github.com/ijalalfrz/coupon-system/internal/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

func TestCouponService_CreateCoupon(t *testing.T) {
	// Dependencies initialized once
	ctx := context.Background()

	// test helper
	createCouponTest := func(
		name string,
		req dto.CreateCouponRequest,
		setupMocks func(*MockCouponRepositoryWriter),
		wantErr bool,
	) func(t *testing.T) {
		return func(t *testing.T) {
			mockWriter := NewMockCouponRepositoryWriter(t)
			setupMocks(mockWriter)

			svc := NewCouponService(mockWriter, nil, nil)
			err := svc.CreateCoupon(ctx, req)

			if wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	}

	t.Run("success", createCouponTest(
		"success",
		dto.CreateCouponRequest{Name: "PROMO10", Amount: 10},
		func(m *MockCouponRepositoryWriter) {
			m.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Return(nil)
		},
		false,
	))

	t.Run("transaction_error", createCouponTest(
		"transaction_error",
		dto.CreateCouponRequest{Name: "PROMO10", Amount: 10},
		func(m *MockCouponRepositoryWriter) {
			m.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Return(errors.New("db error"))
		},
		true,
	))

	t.Run("already_exists", createCouponTest(
		"already_exists",
		dto.CreateCouponRequest{Name: "PROMO10", Amount: 10},
		func(m *MockCouponRepositoryWriter) {
			m.EXPECT().
				WithTransaction(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
					// Simulate the transaction function being called
					_ = txFunc(ctx, nil)
				}).
				Return(errors.New("record already exists"))

			m.EXPECT().
				CreateTx(mock.Anything, mock.Anything, mock.Anything).
				Return(errors.New("record already exists"))
		},
		true,
	))
}

func TestCouponService_GetCouponByCouponName(t *testing.T) {
	// Dependencies initialized once
	ctx := context.Background()

	// test helper
	getCouponTest := func(
		name string,
		req dto.GetCouponRequest,
		setupMocks func(*MockCouponRepositoryReader, *MockCouponClaimRepositoryReader),
		want dto.GetCouponResponse,
		wantErr bool,
	) func(t *testing.T) {
		return func(t *testing.T) {
			mockReader := NewMockCouponRepositoryReader(t)
			mockClaimReader := NewMockCouponClaimRepositoryReader(t)
			setupMocks(mockReader, mockClaimReader)

			svc := NewCouponService(nil, mockReader, mockClaimReader)
			got, err := svc.GetCouponByCouponName(ctx, req)

			if wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !wantErr {
				if diff := cmp.Diff(want, got); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}
			}
		}
	}

	t.Run("success", getCouponTest(
		"success",
		dto.GetCouponRequest{Name: "PROMO10"},
		func(mr *MockCouponRepositoryReader, mcr *MockCouponClaimRepositoryReader) {
			mr.EXPECT().
				GetCouponByName(mock.Anything, "PROMO10").
				Return(model.Coupon{
					ID:              "1",
					Name:            "PROMO10",
					Amount:          10,
					RemainingAmount: 8,
				}, nil)
			mcr.EXPECT().
				GetAllCouponClaimByCouponName(mock.Anything, "PROMO10").
				Return([]string{"user1", "user2"}, nil)
		},
		dto.GetCouponResponse{
			Name:            "PROMO10",
			Amount:          10,
			RemainingAmount: 8,
			ClaimedBy:       []string{"user1", "user2"},
		},
		false,
	))

	t.Run("success_no_claims", getCouponTest(
		"success_no_claims",
		dto.GetCouponRequest{Name: "NEWPROMO"},
		func(mr *MockCouponRepositoryReader, mcr *MockCouponClaimRepositoryReader) {
			mr.EXPECT().
				GetCouponByName(mock.Anything, "NEWPROMO").
				Return(model.Coupon{
					ID:              "2",
					Name:            "NEWPROMO",
					Amount:          5,
					RemainingAmount: 5,
				}, nil)
			mcr.EXPECT().
				GetAllCouponClaimByCouponName(mock.Anything, "NEWPROMO").
				Return([]string{}, nil)
		},
		dto.GetCouponResponse{
			Name:            "NEWPROMO",
			Amount:          5,
			RemainingAmount: 5,
			ClaimedBy:       []string{},
		},
		false,
	))

	t.Run("coupon_not_found", getCouponTest(
		"coupon_not_found",
		dto.GetCouponRequest{Name: "NOTFOUND"},
		func(mr *MockCouponRepositoryReader, mcr *MockCouponClaimRepositoryReader) {
			mr.EXPECT().
				GetCouponByName(mock.Anything, "NOTFOUND").
				Return(model.Coupon{}, errors.New("record not found"))
		},
		dto.GetCouponResponse{},
		true,
	))

	t.Run("claim_fetch_error", getCouponTest(
		"claim_fetch_error",
		dto.GetCouponRequest{Name: "PROMO10"},
		func(mr *MockCouponRepositoryReader, mcr *MockCouponClaimRepositoryReader) {
			mr.EXPECT().
				GetCouponByName(mock.Anything, "PROMO10").
				Return(model.Coupon{
					ID:              "1",
					Name:            "PROMO10",
					Amount:          10,
					RemainingAmount: 8,
				}, nil)
			mcr.EXPECT().
				GetAllCouponClaimByCouponName(mock.Anything, "PROMO10").
				Return(nil, errors.New("db error"))
		},
		dto.GetCouponResponse{},
		true,
	))
}

func TestCouponService_CreateCoupon_WithTransaction(t *testing.T) {
	// Test that verifies the transaction callback is properly invoked
	ctx := context.Background()

	t.Run("transaction_callback_invoked", func(t *testing.T) {
		mockWriter := NewMockCouponRepositoryWriter(t)

		// Capture and execute the transaction callback
		mockWriter.EXPECT().
			WithTransaction(mock.Anything, mock.Anything).
			Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
				// Execute the callback to verify it works
				_ = txFunc(ctx, nil)
			}).
			Return(nil)

		mockWriter.EXPECT().
			CreateTx(mock.Anything, mock.Anything, mock.MatchedBy(func(c model.Coupon) bool {
				return c.Name == "PROMO10" && c.Amount == 10 && c.RemainingAmount == 10
			})).
			Return(nil)

		svc := NewCouponService(mockWriter, nil, nil)
		err := svc.CreateCoupon(ctx, dto.CreateCouponRequest{Name: "PROMO10", Amount: 10})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("transaction_callback_error", func(t *testing.T) {
		mockWriter := NewMockCouponRepositoryWriter(t)

		// Capture and execute the transaction callback, return error from CreateTx
		mockWriter.EXPECT().
			WithTransaction(mock.Anything, mock.Anything).
			Run(func(ctx context.Context, txFunc func(context.Context, *sqlx.Tx) error) {
				_ = txFunc(ctx, nil)
			}).
			Return(errors.New("create failed"))

		mockWriter.EXPECT().
			CreateTx(mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("create failed"))

		svc := NewCouponService(mockWriter, nil, nil)
		err := svc.CreateCoupon(ctx, dto.CreateCouponRequest{Name: "PROMO10", Amount: 10})

		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
