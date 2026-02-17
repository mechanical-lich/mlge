package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator is a function that validates a value and returns an error if invalid
type Validator func(value string) error

// Required validates that a value is not empty
func Required(fieldName string) Validator {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return ValidationError{Field: fieldName, Message: "This field is required"}
		}
		return nil
	}
}

// MinLength validates that a value has at least minLen characters
func MinLength(fieldName string, minLen int) Validator {
	return func(value string) error {
		if len(value) < minLen {
			return ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("Must be at least %d characters", minLen),
			}
		}
		return nil
	}
}

// MaxLength validates that a value has at most maxLen characters
func MaxLength(fieldName string, maxLen int) Validator {
	return func(value string) error {
		if len(value) > maxLen {
			return ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("Must be at most %d characters", maxLen),
			}
		}
		return nil
	}
}

// Email validates that a value is a valid email address
func Email(fieldName string) Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return func(value string) error {
		if value != "" && !emailRegex.MatchString(value) {
			return ValidationError{Field: fieldName, Message: "Must be a valid email address"}
		}
		return nil
	}
}

// Numeric validates that a value is a valid number
func Numeric(fieldName string) Validator {
	return func(value string) error {
		if value != "" {
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				return ValidationError{Field: fieldName, Message: "Must be a number"}
			}
		}
		return nil
	}
}

// IntegerRange validates that a value is an integer within a range
func IntegerRange(fieldName string, min, max int) Validator {
	return func(value string) error {
		if value == "" {
			return nil
		}
		num, err := strconv.Atoi(value)
		if err != nil {
			return ValidationError{Field: fieldName, Message: "Must be an integer"}
		}
		if num < min || num > max {
			return ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("Must be between %d and %d", min, max),
			}
		}
		return nil
	}
}

// Pattern validates that a value matches a regex pattern
func Pattern(fieldName, pattern, message string) Validator {
	regex := regexp.MustCompile(pattern)
	return func(value string) error {
		if value != "" && !regex.MatchString(value) {
			return ValidationError{Field: fieldName, Message: message}
		}
		return nil
	}
}

// Custom creates a custom validator
func Custom(fieldName string, validatorFunc func(string) bool, message string) Validator {
	return func(value string) error {
		if !validatorFunc(value) {
			return ValidationError{Field: fieldName, Message: message}
		}
		return nil
	}
}

// Combine combines multiple validators into one
func Combine(validators ...Validator) Validator {
	return func(value string) error {
		for _, validator := range validators {
			if err := validator(value); err != nil {
				return err
			}
		}
		return nil
	}
}
