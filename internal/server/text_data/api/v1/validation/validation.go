package validation

import (
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/go-playground/validator/v10"
)

// Validator struct holds a validator instance for validating requests.
type Validator struct {
	validator *validator.Validate // Instance of the validator for performing validations
}

// New creates a new Validator instance.
func New(validator *validator.Validate) (*Validator, error) {
	v := &Validator{validator: validator}

	return v, nil
}

// ValidatePostRequest validates the incoming request for posting text data.
func (v *Validator) ValidatePostRequest(req *model.TextDataPostRequest) (map[string]string, bool) {
	err := v.validator.Struct(req)
	report := make(map[string]string)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
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
