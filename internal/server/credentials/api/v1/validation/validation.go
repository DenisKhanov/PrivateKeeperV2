package validation

import (
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/go-playground/validator/v10"
)

// Validator struct encapsulates the validator instance from go-playground/validator package.
// It is responsible for validating incoming requests.
type Validator struct {
	validator *validator.Validate // The validator instance for performing validation checks
}

// New initializes a new Validator instance with a validator.
func New(validator *validator.Validate) (*Validator, error) {
	v := &Validator{validator: validator}

	return v, nil
}

// ValidatePostRequest validates the incoming CredentialsPostRequest.
// It checks required fields and returns a map of validation errors if any exist.
func (v *Validator) ValidatePostRequest(req *model.CredentialsPostRequest) (map[string]string, bool) {
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
