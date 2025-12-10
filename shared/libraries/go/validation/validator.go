package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate valide une structure
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidateEmail valide une adresse email
func ValidateEmail(email string) bool {
	err := validate.Var(email, "required,email")
	return err == nil
}

