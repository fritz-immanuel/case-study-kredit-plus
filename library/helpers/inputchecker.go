package helpers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func MultiValueFilterCheck(input string) string {
	result := ""

	if input != "" {
		sanitizeSpace := strings.ReplaceAll(input, " ", "")
		splitComa := strings.Split(sanitizeSpace, ",")

		for _, s := range splitComa {
			num, _ := strconv.Atoi(s)

			if reflect.TypeOf(num).String() == "int" {
				if result == "" {
					result = fmt.Sprintf("%d", num)
				} else {
					result = fmt.Sprintf("%s,%d", result, num)
				}
			}
		}
	}

	return result
}
