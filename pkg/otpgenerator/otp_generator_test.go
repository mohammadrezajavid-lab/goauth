// Package otpgenerator pkg/otpgenerator/otp_generator_test.go
package otpgenerator

import (
	"regexp"
	"strings"
	"testing"
)

func TestNewOTPGenerator(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := Config{OTPChars: "1234567890"}
		generator, err := NewOTPGenerator(config)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if generator == nil {
			t.Error("expected generator to be created")
		}
		if generator.config.OTPChars != "1234567890" {
			t.Errorf("expected OTPChars to be '1234567890', got '%s'", generator.config.OTPChars)
		}
	})

	t.Run("empty OTP chars", func(t *testing.T) {
		config := Config{OTPChars: ""}
		generator, err := NewOTPGenerator(config)

		if err != ErrEmptyChars {
			t.Errorf("expected ErrEmptyChars, got %v", err)
		}
		if generator != nil {
			t.Error("expected generator to be nil")
		}
	})
}

func TestOTPGenerator_GenerateOTP(t *testing.T) {
	config := Config{OTPChars: "1234567890"}
	generator, _ := NewOTPGenerator(config)

	t.Run("valid length", func(t *testing.T) {
		lengths := []int{4, 6, 8, 10}
		for _, length := range lengths {
			t.Run(string(rune('0'+length)), func(t *testing.T) {
				otp, err := generator.GenerateOTP(length)

				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if len(otp) != length {
					t.Errorf("expected OTP length %d, got %d", length, len(otp))
				}

				// Check if all characters are from the configured set
				for _, char := range otp {
					if !strings.ContainsRune(config.OTPChars, char) {
						t.Errorf("OTP contains invalid character: %c", char)
					}
				}
			})
		}
	})

	t.Run("zero length", func(t *testing.T) {
		otp, err := generator.GenerateOTP(0)

		if err != ErrInvalidLength {
			t.Errorf("expected ErrInvalidLength, got %v", err)
		}
		if otp != "" {
			t.Errorf("expected empty string, got %s", otp)
		}
	})

	t.Run("negative length", func(t *testing.T) {
		otp, err := generator.GenerateOTP(-1)

		if err != ErrInvalidLength {
			t.Errorf("expected ErrInvalidLength, got %v", err)
		}
		if otp != "" {
			t.Errorf("expected empty string, got %s", otp)
		}
	})

	t.Run("custom character set", func(t *testing.T) {
		customConfig := Config{OTPChars: "ABCDEF"}
		customGenerator, _ := NewOTPGenerator(customConfig)
		otp, err := customGenerator.GenerateOTP(6)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(otp) != 6 {
			t.Errorf("expected OTP length 6, got %d", len(otp))
		}

		// Check if all characters are from the custom set
		for _, char := range otp {
			if !strings.ContainsRune("ABCDEF", char) {
				t.Errorf("OTP contains invalid character: %c", char)
			}
		}
	})

	t.Run("randomness test", func(t *testing.T) {
		// Generate multiple OTPs and check they're different
		otps := make(map[string]bool)
		for i := 0; i < 100; i++ {
			otp, err := generator.GenerateOTP(6)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			otps[otp] = true
		}

		// We expect at least 80% unique values for 6-digit OTPs
		if len(otps) < 80 {
			t.Errorf("expected at least 80 unique OTPs, got %d", len(otps))
		}
	})
}

func TestGenerateOTP(t *testing.T) {
	t.Run("valid length", func(t *testing.T) {
		lengths := []int{4, 6, 8, 10}
		for _, length := range lengths {
			t.Run(string(rune('0'+length)), func(t *testing.T) {
				otp, err := GenerateOTP(length)

				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if len(otp) != length {
					t.Errorf("expected OTP length %d, got %d", length, len(otp))
				}

				// Check if all characters are numeric
				matched, _ := regexp.MatchString("^[0-9]+$", otp)
				if !matched {
					t.Errorf("OTP should contain only numeric characters, got %s", otp)
				}
			})
		}
	})

	t.Run("zero length", func(t *testing.T) {
		otp, err := GenerateOTP(0)

		if err != ErrInvalidLength {
			t.Errorf("expected ErrInvalidLength, got %v", err)
		}
		if otp != "" {
			t.Errorf("expected empty string, got %s", otp)
		}
	})

	t.Run("negative length", func(t *testing.T) {
		otp, err := GenerateOTP(-1)

		if err != ErrInvalidLength {
			t.Errorf("expected ErrInvalidLength, got %v", err)
		}
		if otp != "" {
			t.Errorf("expected empty string, got %s", otp)
		}
	})

	t.Run("default character set", func(t *testing.T) {
		otp, err := GenerateOTP(10)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Should only contain digits 0-9
		for _, char := range otp {
			if char < '0' || char > '9' {
				t.Errorf("OTP contains invalid character: %c", char)
			}
		}
	})

	t.Run("randomness test", func(t *testing.T) {
		// Generate multiple OTPs and check they're different
		otps := make(map[string]bool)
		for i := 0; i < 100; i++ {
			otp, err := GenerateOTP(6)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			otps[otp] = true
		}

		// We expect at least 80% unique values for 6-digit OTPs
		if len(otps) < 80 {
			t.Errorf("expected at least 80 unique OTPs, got %d", len(otps))
		}
	})
}

// Benchmark tests
func BenchmarkGenerateOTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateOTP(6)
	}
}

func BenchmarkOTPGenerator_GenerateOTP(b *testing.B) {
	config := Config{OTPChars: "1234567890"}
	generator, _ := NewOTPGenerator(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = generator.GenerateOTP(6)
	}
}

// Example test
func ExampleGenerateOTP() {
	otp, err := GenerateOTP(6)
	if err != nil {
		panic(err)
	}
	// OTP will be a 6-digit numeric string like "123456"
	_ = otp
}

func ExampleOTPGenerator_GenerateOTP() {
	config := Config{OTPChars: "ABCDEF123456"}
	generator, err := NewOTPGenerator(config)
	if err != nil {
		panic(err)
	}

	otp, err := generator.GenerateOTP(8)
	if err != nil {
		panic(err)
	}
	// OTP will be an 8-character string using only A-F and 1-6
	_ = otp
}
