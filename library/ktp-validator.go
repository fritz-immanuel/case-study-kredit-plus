package library

import "regexp"

func IsNIKValid(NIK string) bool {
	NIKRegex := regexp.MustCompile(`^\d{16}$`)
	return NIKRegex.MatchString(NIK)
}
