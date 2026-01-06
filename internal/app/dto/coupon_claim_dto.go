package dto

import (
	"fmt"
	"net/http"

	"github.com/ijalalfrz/coupon-system/internal/pkg/exception"
)

type ClaimCouponRequest struct {
	UserID     string `json:"user_id" validate:"required"`
	CouponName string `json:"coupon_name" validate:"required"`
}

func (req *ClaimCouponRequest) Bind(r *http.Request) error {

	err := validate.Struct(req)
	if err != nil {
		return exception.ApplicationError{
			Message:    fmt.Sprintf("validate claim coupon request: %v", err),
			StatusCode: http.StatusBadRequest,
			Cause:      err,
		}
	}

	return nil
}
