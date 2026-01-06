package model

type CouponClaim struct {
	ID         string `json:"id"`
	CouponName string `json:"coupon_name"`
	UserID     string `json:"user_id"`
}
