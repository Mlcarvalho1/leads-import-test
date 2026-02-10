package validation

import (
	"fmt"
	"strconv"

	"github.com/nyaruka/phonenumbers"
)

type PhoneInfo struct {
	DialCode    string
	CountryCode string
	National    string
}

func ParsePhone(raw string) (*PhoneInfo, error) {
	num, err := phonenumbers.Parse(raw, "BR")
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	if !phonenumbers.IsValidNumber(num) {
		return nil, fmt.Errorf("invalid phone number")
	}

	countryCode := phonenumbers.GetRegionCodeForNumber(num)
	dialCode := strconv.Itoa(int(num.GetCountryCode()))
	national := phonenumbers.Format(num, phonenumbers.NATIONAL)

	// Strip non-digits from national number
	cleaned := ""
	for _, c := range national {
		if c >= '0' && c <= '9' {
			cleaned += string(c)
		}
	}

	return &PhoneInfo{
		DialCode:    dialCode,
		CountryCode: countryCode,
		National:    cleaned,
	}, nil
}
