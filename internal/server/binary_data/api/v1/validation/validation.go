package validation

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func New(validator *validator.Validate) (*Validator, error) {
	v := &Validator{validator: validator}

	return v, nil
}
