package validation

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func New(validator *validator.Validate) *Validator {
	return &Validator{validator: validator}
}
