package service

import (
	"net/http"

	"github.com/ijalalfrz/coupon-system/internal/pkg/exception"
)

var ErrCouponStockIsNotEnough = exception.ApplicationError{
	Message:    "coupon stock is not enough",
	StatusCode: http.StatusBadRequest,
}
