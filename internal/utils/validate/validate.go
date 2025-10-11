package validate

import "github.com/go-playground/validator/v10"

func Struct(input any) error {
	validate := validator.New()
	return validate.Struct(input)
}
