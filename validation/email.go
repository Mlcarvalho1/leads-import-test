package validation

import (
	"fmt"
	"net/mail"
)

func ValidateEmail(raw string) error {
	if len(raw) > 255 {
		return fmt.Errorf("email must be at most 255 characters")
	}

	_, err := mail.ParseAddress(raw)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
