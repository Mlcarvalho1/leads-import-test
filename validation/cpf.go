package validation

import (
	"fmt"
	"strings"
)

func ValidateCPF(raw string) (string, error) {
	// Strip non-digits
	digits := ""
	for _, c := range raw {
		if c >= '0' && c <= '9' {
			digits += string(c)
		}
	}

	if len(digits) != 11 {
		return "", fmt.Errorf("CPF must have 11 digits")
	}

	// Check all same digit
	if strings.Count(digits, string(digits[0])) == 11 {
		return "", fmt.Errorf("invalid CPF")
	}

	// First check digit
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	remainder := (sum * 10) % 11
	if remainder == 10 {
		remainder = 0
	}
	if remainder != int(digits[9]-'0') {
		return "", fmt.Errorf("invalid CPF")
	}

	// Second check digit
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	remainder = (sum * 10) % 11
	if remainder == 10 {
		remainder = 0
	}
	if remainder != int(digits[10]-'0') {
		return "", fmt.Errorf("invalid CPF")
	}

	return digits, nil
}
