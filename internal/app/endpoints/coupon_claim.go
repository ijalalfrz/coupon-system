package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/coupon-system/internal/app/dto"
)

type CouponClaimServiceManager interface {
	ClaimCoupon(ctx context.Context, req dto.ClaimCouponRequest) error
}

type CouponClaimEndpoint struct {
	ClaimCoupon endpoint.Endpoint
}

func MakeCouponClaimEndpoint(service CouponClaimServiceManager) CouponClaimEndpoint {
	return CouponClaimEndpoint{
		ClaimCoupon: makeClaimCouponEndpoint(service),
	}
}

func makeClaimCouponEndpoint(service CouponClaimServiceManager) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request, ok := req.(*dto.ClaimCouponRequest)
		if !ok || request == nil {
			return nil, errors.New("invalid type")
		}

		err := service.ClaimCoupon(ctx, *request)
		if err != nil {
			return nil, fmt.Errorf("coupon claim service: %w", err)
		}

		return nil, nil
	}
}
