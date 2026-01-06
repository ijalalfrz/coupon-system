package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/coupon-system/internal/app/dto"
)

type CouponServiceManager interface {
	CreateCoupon(ctx context.Context, req dto.CreateCouponRequest) error
	GetCouponByCouponName(ctx context.Context, req dto.GetCouponRequest) (dto.GetCouponResponse, error)
}

type CouponEndpoint struct {
	CreateCoupon          endpoint.Endpoint
	GetCouponByCouponName endpoint.Endpoint
}

func MakeCouponEndpoint(service CouponServiceManager) CouponEndpoint {
	return CouponEndpoint{
		CreateCoupon:          makeCreateCouponEndpoint(service),
		GetCouponByCouponName: makeGetCouponByCouponNameEndpoint(service),
	}
}

func makeCreateCouponEndpoint(service CouponServiceManager) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request, ok := req.(*dto.CreateCouponRequest)
		if !ok || request == nil {
			return nil, errors.New("invalid type")
		}

		err := service.CreateCoupon(ctx, *request)
		if err != nil {
			return nil, fmt.Errorf("coupon service: %w", err)
		}

		return nil, nil
	}
}

func makeGetCouponByCouponNameEndpoint(service CouponServiceManager) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request, ok := req.(*dto.GetCouponRequest)
		if !ok || request == nil {
			return nil, errors.New("invalid type")
		}

		coupon, err := service.GetCouponByCouponName(ctx, *request)
		if err != nil {
			return nil, fmt.Errorf("coupon service: %w", err)
		}

		return coupon, nil
	}
}
