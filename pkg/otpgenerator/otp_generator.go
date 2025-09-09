// Package otpgenerator pkg/otpgenerator/otp_generator.go
package otpgenerator

import (
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrInvalidLength = errors.New("invalid OTP length: must be greater than 0")
	ErrEmptyChars    = errors.New("OTP characters cannot be empty")
)

type Config struct {
	OTPChars string `koanf:"otp_chars"`
}

type OTPGenerator struct {
	config Config
}

// NewOTPGenerator creates a new OTP generator with the given config
func NewOTPGenerator(config Config) (*OTPGenerator, error) {
	if config.OTPChars == "" {
		return nil, ErrEmptyChars
	}
	return &OTPGenerator{config: config}, nil
}

// GenerateOTP generates an OTP using the configured characters
func (g *OTPGenerator) GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	buffer := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, buffer)
	if err != nil {
		return "", err
	}

	chars := g.config.OTPChars
	for i := 0; i < length; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}

	return string(buffer), nil
}

// GenerateOTP generates an OTP with default numeric characters (standalone function)
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	const defaultOTPChars = "1234567890"
	buffer := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, buffer)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		buffer[i] = defaultOTPChars[int(buffer[i])%len(defaultOTPChars)]
	}

	return string(buffer), nil
}
