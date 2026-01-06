package dto

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ijalalfrz/coupon-system/internal/pkg/exception"
)

type CreateCouponRequest struct {
	Name   string `json:"name" validate:"required"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}

func (req *CreateCouponRequest) Bind(r *http.Request) error {

	err := validate.Struct(req)
	if err != nil {
		return exception.ApplicationError{
			Message:    fmt.Sprintf("validate create coupon request: %v", err),
			StatusCode: http.StatusBadRequest,
			Cause:      err,
		}
	}

	return nil
}

type GetCouponRequest struct {
	Name string `json:"name" validate:"required"`
}

func (req *GetCouponRequest) Bind(r *http.Request) error {
	req.Name = chi.URLParam(r, "name")

	err := validate.Struct(req)
	if err != nil {
		return exception.ApplicationError{
			Message:    fmt.Sprintf("validate get coupon request: %v", err),
			StatusCode: http.StatusBadRequest,
			Cause:      err,
		}
	}

	return nil
}

type GetCouponResponse struct {
	Name            string   `json:"name"`
	Amount          int      `json:"amount"`
	RemainingAmount int      `json:"remaining_amount"`
	ClaimedBy       []string `json:"claimed_by"`
}
