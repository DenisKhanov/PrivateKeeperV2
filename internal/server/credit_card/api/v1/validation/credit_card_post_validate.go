package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (v *Validator) ValidatePostRequest(req *model.CreditCardPostRequest) (map[string]string, bool) {
	err := v.validator.Struct(req)
	report := make(map[string]string)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationErr := range validationErrors {
				switch validationErr.Tag() {
				case "card_number":
					report[validationErr.Field()] = "must be valid card_number"
				case "owner":
					report[validationErr.Field()] = "must be valid owner: first_name second_name"
				case "cvv":
					report[validationErr.Field()] = "must be valid cvv"
				case "pin":
					report[validationErr.Field()] = "must be valid pin"
				case "required":
					report[validationErr.Field()] = "is required"
				case "expires_at":
					report[validationErr.Field()] = "expires_at must be in DD-MM-YYYY format"
				}
			}
			return report, false
		}
		return map[string]string{"error": "unknown validation error"}, false
	}
	return nil, true
}
