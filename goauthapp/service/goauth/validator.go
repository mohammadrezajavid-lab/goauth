package goauth

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

// Validator handles request validation.
type Validator struct {
}

// NewValidator creates a new Validator.
func NewValidator() Validator {
	return Validator{}
}

func (v *Validator) ValidateGenerateOTPRequest(request *GenerateOTPRequest) (map[string]string, error) {
	if err := v.validateGenerateOTPRequest(request); err != nil {
		fieldErrors := make(map[string]string)
		var valueErr validation.Errors
		ok := errors.As(err, &valueErr)
		if ok {
			for key, value := range valueErr {
				fieldErrors[key] = value.Error()
			}
		}

		return fieldErrors, errors.Join(errors.New("validation GenerateOTPRequest error"), err)
	}

	return nil, nil
}

func (v *Validator) ValidateVerifyOTPRequest(request *VerifyOTPRequest) (map[string]string, error) {
	if err := v.validateVerifyOTPRequest(request); err != nil {
		fieldErrors := make(map[string]string)
		var valueErr validation.Errors
		ok := errors.As(err, &valueErr)
		if ok {
			for key, value := range valueErr {
				fieldErrors[key] = value.Error()
			}
		}

		return fieldErrors, errors.Join(errors.New("validation VerifyOTPRequest error"), err)
	}

	return nil, nil
}

// ValidateGenerateOTPRequest validates the request for generating an OTP.
func (v *Validator) validateGenerateOTPRequest(request *GenerateOTPRequest) error {
	return validation.ValidateStruct(request,
		// phoneNumber must be required and in E.164 format (e.g., +989123456789)
		validation.Field(&request.PhoneNumber,
			validation.Required.Error("phone number is required"),
			validation.Match(regexp.MustCompile(`^\+[1-9]\d{1,14}$`)).Error("phone number must be in E.164 format (e.g., +989123456789)"),
		),
	)
}

// ValidateVerifyOTPRequest validates the request for verifying an OTP and logging in.
func (v *Validator) validateVerifyOTPRequest(request *VerifyOTPRequest) error {
	return validation.ValidateStruct(request,
		// phoneNumber must be required and in E.164 format
		validation.Field(&request.PhoneNumber,
			validation.Required.Error("phone number is required"),
			validation.Match(regexp.MustCompile(`^\+[1-9]\d{1,14}$`)).Error("phone number must be in E.164 format"),
		),
		// OTP must be required and exactly 6 digits long
		validation.Field(&request.OTP,
			validation.Required.Error("otp is required"),
			validation.Length(6, 6).Error("otp must be 6 digits"),
		),
	)
}
