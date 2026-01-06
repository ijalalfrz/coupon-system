package repository

import (
	"net/http"

	"github.com/ijalalfrz/coupon-system/internal/pkg/exception"
)

var ErrCouponAlreadyClaimed = exception.ApplicationError{
	Message:    "coupon already claimed by this user",
	StatusCode: http.StatusConflict,
}

var ErrRecordNotFound = exception.ApplicationError{
	Message:    "record not found",
	StatusCode: http.StatusNotFound,
}

var ErrRecordAlreadyExists = exception.ApplicationError{
	Message:    "record already exists",
	StatusCode: http.StatusConflict,
}
