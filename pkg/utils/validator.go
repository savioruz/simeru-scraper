package utils

import "github.com/go-playground/validator/v10"

var daysOfWeek = map[string]bool{
	"senin":  true,
	"selasa": true,
	"rabu":   true,
	"kamis":  true,
	"jumat":  true,
	"sabtu":  true,
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()

	err := validate.RegisterValidation("day", isValidDayOfWeek)
	if err != nil {
		return nil
	}

	return &Validator{validate}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func isValidDayOfWeek(fl validator.FieldLevel) bool {
	day := fl.Field().String()
	return daysOfWeek[day]
}
