package context

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	inner *validator.Validate
}

func NewValidator() Validator {
	return Validator{inner: validator.New()}
}

func (v Validator) Validate(i interface{}) error {
	if err := v.inner.Struct(i); err != nil {
		return err
	}
	return nil
}
