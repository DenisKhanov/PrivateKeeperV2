package validation

import (
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New(validator *validator.Validate) (*Validator, error) {
	v := &Validator{validator: validator}

	return v, nil
}

func (v *Validator) ValidatePostRequest(req *model.TextDataPostRequest) (map[string]string, bool) {
	err := v.validator.Struct(req)
	report := make(map[string]string)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() { //nolint:gocritic
				case "required":
					report[validationErr.Field()] = "is required"
				}
			}
			return report, false
		}
		return map[string]string{"error": "unknown validation error"}, false
	}
	return nil, true
}
