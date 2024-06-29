package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

type Validator = *validator.Validate

func GetValidatorInstance() Validator {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())
	}
	return validate
}
