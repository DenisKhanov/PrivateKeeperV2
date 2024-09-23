package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

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
