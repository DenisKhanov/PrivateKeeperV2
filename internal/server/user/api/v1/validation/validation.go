package validation

import (
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/go-playground/validator/v10"
)

// Validator struct encapsulates a validator instance
type Validator struct {
	validator *validator.Validate // Validator instance to perform validation
}

// New initializes a new Validator instance
func New(validator *validator.Validate) *Validator {
	return &Validator{validator: validator}
}

// ValidateLoginRequest validates the login request structure
func (v *Validator) ValidateLoginRequest(req *model.UserLoginRequest) (map[string]string, bool) {
	err := v.validator.Struct(req)
	report := make(map[string]string)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
				case "email":
					report[validationErr.Field()] = "must be valid email"
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

// ValidateRegisterRequest validates the registration request structure
func (v *Validator) ValidateRegisterRequest(req *model.UserRegisterRequest) (map[string]string, bool) {
	err := v.validator.Struct(req)
	report := make(map[string]string)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
				case "email":
					report[validationErr.Field()] = "must be valid email"
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
