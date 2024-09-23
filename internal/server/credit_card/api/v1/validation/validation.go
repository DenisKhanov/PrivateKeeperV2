package validation

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const expiresAtLayout = "02-01-2006"

type Validator struct {
	validator *validator.Validate
}

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

func expiresAt(fl validator.FieldLevel) bool {
	_, err := time.Parse(expiresAtLayout, fl.Field().String())
	return err == nil
}

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

func owner(fl validator.FieldLevel) bool {
	return len(strings.Split(fl.Field().String(), " ")) == 2
}
