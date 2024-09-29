package validation

import (
	"errors"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const expiresAtLayout = "02-01-2006" // Define the layout for the expiration date format

// Validator is a struct that holds the validator instance.
type Validator struct {
	validator *validator.Validate // Instance of the validator
}

// New creates a new Validator and registers custom validation rules.
func New(validator *validator.Validate) (*Validator, error) {
	v := &Validator{validator: validator}

	err := v.validator.RegisterValidation("cvv", cvvCode)
	if err != nil {
		return nil, fmt.Errorf("register cvv: %w", err)
	}

	err = v.validator.RegisterValidation("pin", pinCode)
	if err != nil {
		return nil, fmt.Errorf("register pin: %w", err)
	}

	err = v.validator.RegisterValidation("owner", owner)
	if err != nil {
		return nil, fmt.Errorf("register owner: %w", err)
	}

	err = v.validator.RegisterValidation("card_number", cardNumber)
	if err != nil {
		return nil, fmt.Errorf("register card_number: %w", err)
	}

	err = v.validator.RegisterValidation("expires_at", expiresAt)
	if err != nil {
		return nil, fmt.Errorf("register expires_at: %w", err)
	}

	return v, nil
}

// expiresAt checks if the expiration date is in the correct format (DD-MM-YYYY).
func expiresAt(fl validator.FieldLevel) bool {
	_, err := time.Parse(expiresAtLayout, fl.Field().String())
	return err == nil
}

// cardNumber validates the card number format (four groups of four digits).
func cardNumber(fl validator.FieldLevel) bool {
	blocks := strings.Split(fl.Field().String(), " ")
	if len(blocks) != 4 {
		return false
	}

	for _, block := range blocks {
		if len(block) != 4 {
			return false
		}
		for _, char := range block {
			if unicode.IsLetter(char) {
				return false
			}
		}
	}
	return true
}

// cvvCode validates the CVV code (must be three digits).
func cvvCode(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) != 3 {
		return false
	}
	for _, char := range fl.Field().String() {
		if unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// pinCode validates the PIN code (must be four digits).
func pinCode(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) != 4 {
		return false
	}
	for _, char := range fl.Field().String() {
		if unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// owner validates the owner's name (must consist of two words).
func owner(fl validator.FieldLevel) bool {
	return len(strings.Split(fl.Field().String(), " ")) == 2
}

// ValidatePostRequest validates the incoming CreditCardPostRequest.
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
