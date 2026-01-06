package model

type Coupon struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Amount          int    `json:"amount"`
	RemainingAmount int    `json:"remaining_amount"`
}
