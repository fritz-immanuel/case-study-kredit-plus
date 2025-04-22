package library

import "regexp"

func ValidateCountryCode(code string) bool {
	var countryCodePattern = `^\+[1-9]\d{0,2}$`

	re := regexp.MustCompile(countryCodePattern)

	return re.MatchString(code)
}

func ValidatePhoneNumber(phone string) bool {
	var phonePattern = `^[1-9]\d{1,14}$`

	re := regexp.MustCompile(phonePattern)

	return re.MatchString(phone)
}
