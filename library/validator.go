package library

import (
	"regexp"

	"github.com/google/uuid"
)

var validTenors = map[int]struct{}{
	1: {}, 2: {}, 3: {}, 6: {},
}

func ValidateTenor(tenor int) bool {
	_, ok := validTenors[tenor]

	return ok
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateNIK(NIK string) bool {
	NIKRegex := regexp.MustCompile(`^\d{16}$`)
	return NIKRegex.MatchString(NIK)
}

func ValidateCountryCode(code string) bool {
	countryCodeRegex := regexp.MustCompile(`^\+[1-9]\d{0,2}$`)
	return countryCodeRegex.MatchString(code)
}

func ValidatePhoneNumber(phone string) bool {
	phoneNumberRegex := regexp.MustCompile(`^[1-9]\d{1,14}$`)
	return phoneNumberRegex.MatchString(phone)
}

func ValidateTextInput(input string) bool {
	inputRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-\',.&()]{1,250}$`)
	return inputRegex.MatchString(input)
}

func ValidateUUID(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil
}
