package dto

import "github.com/go-playground/validator/v10"

var (
	validate = validator.New()
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct {
	Message string `json:"message"`
}
